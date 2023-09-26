package operations

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"

	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"

	"github.com/sabahtalateh/oapi/internal"
	"github.com/sabahtalateh/oapi/internal/node/maps"
)

type Path struct {
	method    string
	url       string
	normUrl   string
	params    []param
	responses []Response
	requests  []Request
}

type param struct {
	name     string
	in       string
	required bool
	typ      string
}

func (p Path) operation() {}

func (p Path) Sync(ctx internal.Context, n *yaml.Node) {
	logSyncPath(ctx, p)

	paths := maps.MakePath(n, "paths")
	p.sync(paths)
}

func parsePath(ctx internal.Context, args []string) (Path, error) {
	if len(args) < 2 {
		return Path{}, errors.Join(errMalformedPath, errPathUsage)
	}

	path := Path{}

	path.method = args[0]
	if !slices.Contains(methods, path.method) {
		return path, errors.Join(
			fmt.Errorf("unknown method: %s\nexpected one of: %s", path.method, strings.Join(methods, ", ")),
			errPathUsage,
		)
	}

	path.method = strings.ToLower(path.method)

	path.url = args[1]
	path.params = pathParams(path.url)
	path.normUrl = path.url
	for _, p := range path.params {
		path.normUrl = strings.ReplaceAll(path.normUrl, p.name, "")
	}

	pathFunc, fset, err := pathFunc(ctx)
	if err != nil {
		return path, err
	}

	path.requests, err = collectRequests(ctx.WithFSet(fset), pathFunc)
	if err != nil {
		return path, err
	}

	path.responses, err = collectResponses(ctx.WithFSet(fset), pathFunc)
	if err != nil {
		return path, err
	}

	return path, nil
}

func pathParams(url string) []param {
	inParam := false
	var parameter string
	var params []param
	for _, x := range url {
		if x == '{' {
			inParam = true
			continue
		}

		if x == '}' {
			params = append(params, param{name: parameter, in: "path", required: true, typ: "string"})
			parameter = ""
			inParam = false
		}

		if inParam {
			parameter += string(x)
		}
	}

	return params
}

func pathFunc(ctx internal.Context) (*ast.FuncDecl, *token.FileSet, error) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, ctx.Location.File, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}

	var (
		fun   *ast.FuncDecl
		gDecl *ast.GenDecl
	)

	ast.Inspect(file, func(n ast.Node) bool {
		if fun != nil || gDecl != nil {
			return false
		}
		if n == nil {
			return true
		}
		switch nn := n.(type) {
		case *ast.GenDecl:
			if fset.Position(n.Pos()).Line > ctx.Location.Line {
				gDecl = nn
				return false
			}
		case *ast.FuncDecl:
			if fset.Position(n.Pos()).Line > ctx.Location.Line {
				fun = nn
				return false
			}
		}

		return true
	})

	err = errors.Join(errMalformedPath, e("path should be defined over func"), errPathUsage)
	if gDecl != nil {
		return nil, nil, err
	}
	if fun == nil {
		return nil, nil, err
	}

	return fun, fset, nil
}

func collectRequests(ctx internal.Context, f *ast.FuncDecl) ([]Request, error) {
	if f.Doc == nil {
		return nil, nil
	}

	var (
		out []Request
	)

	for _, comment := range f.Doc.List {
		ff := strings.Fields(comment.Text)
		if len(ff) < 1 {
			continue
		}
		if !strings.HasSuffix(ff[0], "go:generate") {
			continue
		}
		if len(ff) < 2 {
			continue
		}
		if ff[1] != ctx.Executable {
			continue
		}
		if len(ff) < 3 {
			continue
		}
		if ff[2] != "request" {
			continue
		}

		commentPos := ctx.FSet.Position(comment.Pos())
		loc := internal.Location{
			File: commentPos.Filename,
			Line: commentPos.Line,
		}
		request, err := parseRequest(ff[3:])
		checkLoc(loc, err)

		out = append(out, request)
	}

	return out, nil
}

func collectResponses(ctx internal.Context, f *ast.FuncDecl) ([]Response, error) {
	if f.Doc == nil {
		return nil, nil
	}

	var (
		out []Response
	)

	for _, comment := range f.Doc.List {
		ff := strings.Fields(comment.Text)
		if len(ff) < 1 {
			continue
		}
		if !strings.HasSuffix(ff[0], "go:generate") {
			continue
		}
		if len(ff) < 2 {
			continue
		}
		if ff[1] != ctx.Executable {
			continue
		}
		if len(ff) < 3 {
			continue
		}
		if ff[2] != "response" {
			continue
		}

		commentPos := ctx.FSet.Position(comment.Pos())
		loc := internal.Location{
			File: commentPos.Filename,
			Line: commentPos.Line,
		}
		response, err := parseResponse(ff[3:])
		checkLoc(loc, err)

		out = append(out, response)
	}

	return out, nil
}

func checkLoc(l internal.Location, err error) {
	if err != nil {
		fmt.Printf("oapi\n%s\n\t%s:%d\n", err, l.File, l.Line)
		os.Exit(1)
	}
}
