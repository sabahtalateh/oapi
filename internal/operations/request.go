package operations

import (
	"errors"
	"fmt"

	"gopkg.in/yaml.v3"

	"github.com/sabahtalateh/oapi/internal"
)

type Request struct {
	contentType string
	schema      string
	ref         string
}

func (r Request) operation() {}

// Sync do nothing. Synced within Path.Sync
func (r Request) Sync(ctx internal.Context, _ *yaml.Node) {
	logSyncRequest(ctx, r)
	return
}

func parseRequest(args []string) (Request, error) {
	if len(args) < 1 {
		return Request{}, errors.Join(errMalformedRequest, errRequestUsage)
	}

	contentType := ""
	schema := ""

	if len(args) == 1 {
		contentType = "application/json"
		schema = args[0]
	} else {
		contentType = args[0]
		schema = args[1]
	}

	ref := fmt.Sprintf("#/components/schemas/%s", schema)

	return Request{contentType: contentType, schema: schema, ref: ref}, nil
}
