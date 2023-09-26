package operations

import (
	"errors"
	"fmt"
	"go/ast"
	"strings"
	"unicode"

	"github.com/sabahtalateh/oapi/internal"
)

// 1'st return value denotes if it is pointer or not
// 2'nd return value can be of *ast.Ident, *ast.StarExpr or *ast.ArrayType type
func fieldType(t ast.Expr, filedName string) (bool, ast.Expr, error) {
	switch typ := t.(type) {
	case *ast.Ident, *ast.SelectorExpr, *ast.ArrayType:
		return false, typ, nil
	case *ast.StarExpr:
		switch x := typ.X.(type) {
		case *ast.Ident, *ast.SelectorExpr, *ast.ArrayType:
			return true, x, nil
		default:
			return false, nil, fmt.Errorf("unsupported type for property: %s", filedName)
		}
	case *ast.StructType:
		return false, nil, fmt.Errorf("nested struct not supported: %s", filedName)
	default:
		return false, nil, fmt.Errorf("unsupported type for property: %s", filedName)
	}
}

func parseFieldType(ctx internal.Context, typ ast.Expr, propName string) (OASchemaVal, error) {
	var err error

	switch t := typ.(type) {
	case *ast.Ident:
		if _, ok := builtinTypes[t.Name]; ok {
			if t.Name == "byte" {
				return Type{Type: "string", Format: "byte"}, nil
			}

			var pVal Type
			pVal.Type, pVal.Format = builtinToType(t.Name)
			return pVal, nil
		} else {
			var pVal Ref
			pVal.Ref, err = typeToRef(ctx, propName, t)
			if err != nil {
				return nil, err
			}
			return pVal, nil
		}
	case *ast.SelectorExpr:
		if isStdTime(t, ctx.Imports) {
			return Type{Type: "string", Format: "date-time"}, nil
		} else {
			var pVal Ref
			pVal.Ref, err = typeToRef(ctx, propName, t)
			if err != nil {
				return nil, err
			}
			return pVal, nil
		}
	case *ast.ArrayType:
		switch et := t.Elt.(type) {
		case *ast.Ident:
			if et.Name == "byte" {
				return Type{Type: "string", Format: "byte"}, nil
			}
		}
		return parseArrayType(ctx, t)
	default:
		return nil, fmt.Errorf("unsupported type for property: %s", propName)
	}
}

func checkField(ctx internal.Context, f *ast.Field) error {
	if f.Names != nil {
		if len(f.Names) > 1 {
			var names []string
			for _, n := range f.Names {
				names = append(names, n.Name)
			}
			loc := ctx.FSet.Position(f.Names[0].Pos())
			return errors.Join(
				errMalformedSchema,
				fmt.Errorf(
					"multiple fields syntax not allowed. separated fields: %s\n\t%s:%d",
					strings.Join(names, ", "),
					loc.Filename,
					loc.Line,
				),
			)
		}
	}

	return nil
}

// first return - is embedded field
// second return - is field serializable
func fieldInfo(f *ast.Field) (bool, bool, string) {
	if f.Tag != nil {
		tags := strings.Fields(strings.Trim(f.Tag.Value, "`"))
		for _, tag := range tags {
			if strings.HasPrefix(tag, "json:") {
				fieldName := strings.TrimPrefix(tag, "json:")
				fieldName = strings.Trim(fieldName, "\"")
				fieldName = strings.Trim(fieldName, "'")
				return false, true, fieldName
			}
		}
	}

	if f.Names != nil && len(f.Names) != 0 {
		name := f.Names[0].Name
		if name != "" {
			return false, unicode.IsUpper(rune(name[0])), name
		}
	}

	serializable, fieldName := exprInfo(f.Type)
	return true, serializable, fieldName
}

func exprInfo(e ast.Expr) (bool, string) {
	switch typ := e.(type) {
	case *ast.Ident:
		return identInfo(typ)
	case *ast.SelectorExpr:
		return selectorExprInfo(typ)
	case *ast.ArrayType:
		return arrayTypeInfo(typ)
	case *ast.StarExpr:
		return starExprInfo(typ)
	default:
		return false, ""
	}
}

func identInfo(i *ast.Ident) (bool, string) {
	return i.Name != "" && unicode.IsUpper(rune(i.Name[0])), i.Name
}

func selectorExprInfo(s *ast.SelectorExpr) (bool, string) {
	return identInfo(s.Sel)
}

func arrayTypeInfo(a *ast.ArrayType) (bool, string) {
	return exprInfo(a.Elt)
}

func starExprInfo(s *ast.StarExpr) (bool, string) {
	return exprInfo(s.X)
}
