package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema SchemaA
type A struct {
	Field1 string `json:"field1"`
	B      B      `json:"b"`
}

//go:generate oapi schema SchemaB
type B struct {
	C
	Field2 string `json:"field2"`
}

//go:generate oapi schema SchemaC
type C struct {
	Field3 int `json:"field3"`
}

var original = tests.TrimOApi(`
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
          properties:
            field1:
              abc: def
              type: string
    SchemaB:
      allOf:
        - $ref: '#/components/schemas/SchemaC'
    SchemaC:
      type: object
      properties:
        field2:
          type: integer
          format: int64
`)

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
        - b
      properties:
        field1:
          type: string
          abc: def
        b:
          $ref: '#/components/schemas/SchemaB'
    SchemaB:
      allOf:
        - $ref: '#/components/schemas/SchemaC'
        - type: object
          required:
            - field2
          properties:
            field2:
              type: string
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
	tests.WriteOriginal(t, original)

	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	oapi, _ := tests.ReadOAPI(t)
	require.Equal(t, expected, oapi)
}
