package operations

import (
	"errors"

	"gopkg.in/yaml.v3"

	"github.com/sabahtalateh/oapi/internal"
)

var methods = []string{
	"GET",
	"HEAD",
	"POST",
	"PUT",
	"DELETE",
	"CONNECT",
	"OPTIONS",
	"TRACE",
}

type Operation interface {
	operation()
	Sync(ctx internal.Context, n *yaml.Node)
}

func Parse(ctx internal.Context, args []string) (Operation, error) {
	if len(args) == 0 {
		return nil, errors.Join(errUnknownOperation, errUsage)
	}

	token1 := args[0]
	switch token1 {
	case "path":
		return parsePath(ctx, args[1:])
	case "schema":
		return parseSchema(ctx, args[1:])
	case "request":
		return parseRequest(args[1:])
	case "response":
		return parseResponse(args[1:])
	default:
		return nil, errors.Join(errUnknownOperation, e(token1), errUsage)
	}
}
