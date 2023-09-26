package path_added

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema SchemaA
type A struct {
	Field1  string  `json:"field1"`
	CC      []C     `json:"cc"`
	Numbers []int32 `json:"numbers"`
	private string
	Public  string
}

//go:generate oapi schema SchemaAA
type AA []A

//go:generate oapi schema SchemaC
type C struct {
	Field3 int `json:"field3"`
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
        - field1
        - cc
        - numbers
        - Public
      properties:
        field1:
          type: string
        cc:
          type: array
          items:
            $ref: '#/components/schemas/SchemaC'
        numbers:
          type: array
          items:
            type: integer
            format: int32
        Public:
          type: string
    SchemaAA:
      type: array
      items:
        $ref: '#/components/schemas/SchemaA'
    SchemaC:
      type: object
      required:
        - field3
      properties:
        field3:
          type: integer
          format: int64
`)

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	oapi, _ := tests.ReadOAPI(t)
	require.Equal(t, expected, oapi)
}
