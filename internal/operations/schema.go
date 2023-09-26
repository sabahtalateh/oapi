package operations

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"

	"gopkg.in/yaml.v3"

	"github.com/sabahtalateh/oapi/internal"
	"github.com/sabahtalateh/oapi/internal/node/maps"
)

type Schema struct {
	Name   string
	Schema OASchema
}

func (s Schema) operation() {}

func (s Schema) Sync(ctx internal.Context, n *yaml.Node) {
	logSyncSchema(ctx, s)
	schemasN := maps.MakePath(n, "components", "schemas")
	s.sync(schemasN)
}

func parseSchema(ctx internal.Context, args []string) (Schema, error) {
	if len(args) != 1 {
		return Schema{}, errors.Join(errMalformedSchema, errSchemaUsage)
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, ctx.Location.File, nil, parser.ParseComments)
	if err != nil {
		return Schema{}, err
	}

	var typeSpec *ast.TypeSpec
	ast.Inspect(file, func(n ast.Node) bool {
		if typeSpec != nil {
			return false
		}
		if n == nil {
			return true
		}
		switch typ := n.(type) {
		case *ast.TypeSpec:
			if fset.Position(n.Pos()).Line > ctx.Location.Line {
				typeSpec = typ
				return false
			}
		}

		return true
	})

	if typeSpec == nil {
		return Schema{}, errors.Join(errMalformedSchema, errSchemaUsage)
	}

	s, err := parseSchemaType(ctx.WithImports(file.Imports).WithFSet(fset), typeSpec)
	if err != nil {
		return Schema{}, err
	}

	return Schema{Name: args[0], Schema: s}, nil
}
