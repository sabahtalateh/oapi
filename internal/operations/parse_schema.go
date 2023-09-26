package operations

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/sabahtalateh/oapi/internal"
)

type OASchema interface {
	schema()
}

type Object struct {
	Required   []string
	Properties []Property
	Refs       []string
}

// Array can be both OASchema and OASchemaVal
type Array struct {
	Items OASchemaVal // Type or Ref
}

func (o Object) schema() {}
func (a Array) schema()  {}
func (t Type) schema()   {}

type Property struct {
	Name string
	Val  OASchemaVal
}

type OASchemaVal interface {
	schemaVal()
}

type Type struct {
	Type   string
	Format string
}

type Ref struct {
	Ref string
}

func (t Type) schemaVal()  {}
func (r Ref) schemaVal()   {}
func (a Array) schemaVal() {}

func parseSchemaType(ctx internal.Context, spec *ast.TypeSpec) (OASchema, error) {
	switch typ := spec.Type.(type) {
	case *ast.StructType:
		return parseStructType(ctx.WithSchemaType(spec.Name.Name), typ)
	case *ast.ArrayType:
		switch elt := typ.Elt.(type) {
		case *ast.Ident:
			if elt.Name == "byte" {
				return Type{Type: "string", Format: "byte"}, nil
			}
		}
		return parseArrayType(ctx, typ)
	default:
		return Object{}, errors.Join(
			errMalformedSchema,
			e("schema should be defined over struct, array or slice. type defs and type aliases not supported"),
			errSchemaUsage,
		)
	}
}

func parseStructType(ctx internal.Context, Struct *ast.StructType) (Object, error) {
	var (
		o   Object
		err error
	)

	if Struct.Fields == nil {
		return o, nil
	}

	for _, field := range Struct.Fields.List {
		var (
			prop Property
		)

		if err = checkField(ctx, field); err != nil {
			return o, err
		}

		embedded, serializable, fieldName := fieldInfo(field)
		if !serializable {
			continue
		}

		pointer, typ, err := fieldType(field.Type, fieldName)
		if err != nil {
			return o, err
		}

		if embedded {
			ref, err := typeToRef(ctx, fieldName, typ)
			if err != nil {
				return o, err
			}
			o.Refs = append(o.Refs, ref)
			continue
		}

		prop.Name = fieldName

		prop.Val, err = parseFieldType(ctx, typ, prop.Name)
		if err != nil {
			return o, err
		}

		if !pointer {
			o.Required = append(o.Required, prop.Name)
		}

		o.Properties = append(o.Properties, prop)
	}

	return o, nil
}

func parseArrayType(ctx internal.Context, array *ast.ArrayType) (Array, error) {
	val, err := parseFieldType(ctx, array.Elt, "array")
	if err != nil {
		return Array{}, err
	}
	return Array{Items: val}, nil
}

func typeToRef(ctx internal.Context, fieldName string, t ast.Expr) (string, error) {
	switch tt := t.(type) {
	case *ast.Ident:
		return schemaNameForType(ctx, ctx.WorkDir, tt.Name)
	case *ast.SelectorExpr:
		switch alias := tt.X.(type) {
		case *ast.Ident:
			imp, err := findImportPath(ctx, fieldName, alias.Name)
			if err != nil {
				return "", err
			}
			pkgPath, err := internal.PkgDirForImport(ctx.Location, imp)
			if err != nil {
				return "", err
			}
			return schemaNameForType(ctx, pkgPath, tt.Sel.Name)
		default:
			return "", fmt.Errorf("unsupported type for property: %s", fieldName)
		}
	default:
		return "", fmt.Errorf("unsupported type for property: %s", fieldName)
	}
}

func findImportPath(ctx internal.Context, fieldName string, alias string) (string, error) {
	// suppose package dir equals to last package element (dir: c; package a/b/c)
	for _, imp := range ctx.Imports {
		val := impV(imp.Path.Value)

		if imp.Name != nil && imp.Name.Name == alias {
			return val, nil
		}

		if filepath.Base(val) == alias {
			return val, nil
		}
	}

	// if dir not equals to last package element (dir d; package a/b/c)
	// then read ".go" file from dir to read actual package
	for _, imp := range ctx.Imports {
		val := impV(imp.Path.Value)

		readAlias, err := internal.ReadPackageName(ctx.Location, val)
		if err != nil {
			return "", err
		}

		if readAlias == alias {
			return val, nil
		}
	}
	return "", fmt.Errorf("unsupported type for property: %s", fieldName)
}

func schemaNameForType(ctx internal.Context, path string, Type string) (string, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return "", err
	}

	var (
		decl *ast.GenDecl
		typ  *ast.TypeSpec
		cm   ast.CommentMap
	)

	for _, pkg := range pkgs {
		for _, f := range pkg.Files {
			ast.Inspect(f, func(node ast.Node) bool {
				if decl != nil {
					return false
				}
				switch d := node.(type) {
				case *ast.GenDecl:
					for _, spec := range d.Specs {
						switch t := spec.(type) {
						case *ast.TypeSpec:
							if t.Name.Name == Type {
								decl = d
								typ = t
								cm = ast.NewCommentMap(fset, f, f.Comments)
								return false
							}
						}
					}
				}
				return true
			})
		}
	}

	if decl == nil {
		return "", fmt.Errorf("type not found: %s.%s", path, Type)
	}

	if c, ok := cm[decl]; ok {
		for _, group := range c {
			for _, comment := range group.List {
				ref, ok := schemaRefFromComment(ctx, comment.Text)
				if ok {
					return ref, nil
				}
				continue
			}
		}
	}

	loc := fset.Position(typ.Pos())
	return "", errors.Join(
		fmt.Errorf("type %s\n\t%s:%d", typ.Name, loc.Filename, loc.Line),
		fmt.Errorf("can not be used as part of type %s\nas it is not marked with `//go:generate oapi schema ..`", ctx.CurrentSchema),
	)
}

func schemaRefFromComment(ctx internal.Context, comment string) (string, bool) {
	parts := strings.Fields(comment)
	if len(parts) < 4 {
		return "", false
	}
	if !strings.HasSuffix(parts[0], "go:generate") {
		return "", false
	}
	if parts[1] != ctx.Executable {
		return "", false
	}
	if parts[2] != "schema" {
		return "", false
	}
	return fmt.Sprintf("#/components/schemas/%s", parts[3]), true
}

func isStdTime(sel *ast.SelectorExpr, imps []*ast.ImportSpec) bool {
	switch sel.X.(type) {
	case *ast.Ident:
		for _, imp := range imps {
			val := impV(imp.Path.Value)
			if val == "time" && sel.Sel.Name == "Time" {
				return true
			}
		}
	}
	return false
}

func impV(imp string) string {
	v := strings.Trim(imp, "\"")
	v = strings.Trim(v, "'")
	v = strings.Trim(v, "`")

	return v
}
