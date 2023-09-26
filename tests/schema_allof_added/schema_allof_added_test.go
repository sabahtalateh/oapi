package path_added

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema SchemaA
type A struct {
	Field1 string `json:"field1"`
	B
}

//go:generate oapi schema SchemaB
type B struct {
	C
}

//go:generate oapi schema SchemaC
type C struct {
	Field2 int `json:"field2"`
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
      allOf:
        - $ref: '#/components/schemas/SchemaB'
        - type: object
          required:
            - field1
          properties:
            field1:
              type: string
    SchemaB:
      allOf:
        - $ref: '#/components/schemas/SchemaC'
    SchemaC:
      type: object
      required:
        - field2
      properties:
        field2:
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
