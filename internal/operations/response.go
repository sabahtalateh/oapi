package operations

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/sabahtalateh/oapi/internal"
)

type Response struct {
	response    string
	contentType string
	schema      string
	ref         string
}

func (r Response) operation() {}

// Sync do nothing. Synced within Path.Sync
func (r Response) Sync(ctx internal.Context, _ *yaml.Node) {
	logSyncResponse(ctx, r)
	return
}

func parseResponse(args []string) (Response, error) {
	if len(args) < 2 {
		return Response{}, errors.Join(errMalformedResponse, errResponseUsage)
	}

	resp := args[0]
	contentType := ""
	schema := ""

	if len(args) == 2 {
		contentType = "application/json"
		schema = args[1]
	} else {
		contentType = args[1]
		schema = args[2]
	}

	ref := fmt.Sprintf("#/components/schemas/%s", schema)

	return Response{response: resp, contentType: contentType, schema: schema, ref: ref}, nil
}
