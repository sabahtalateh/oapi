package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema SchemaA
type A struct {
	Bytes []byte `json:"bytes"`
}

var expected = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
components:
  schemas:
    SchemaA:
      type: object
      required:
        - bytes
      properties:
        bytes:
          type: string
          format: byte
`)

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	oapi, _ := tests.ReadOAPI(t)
	require.Equal(t, expected, oapi)
}
