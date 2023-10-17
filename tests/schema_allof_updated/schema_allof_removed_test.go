package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema SchemaA
type A struct {
	C
	Field1 string `json:"field1"`
}

//go:generate oapi schema SchemaB
type B struct {
	Field2 string `json:"field2"`
}

//go:generate oapi schema SchemaС
type C struct {
	Field3 string `json:"field3"`
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
        - type: object
          properties:
            field1:
              type: string
              example: Hello1
        - type: object
          properties:
            field1:
              type: string
              example: Hello2
        - $ref: '#/components/schemas/SchemaB'
    SchemaB:
      type: object
      properties:
        field2:
          type: string
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
      allOf:
        - $ref: '#/components/schemas/SchemaС'
        - type: object
          required:
            - field1
          properties:
            field1:
              type: string
              example: Hello1
    SchemaB:
      type: object
      required:
        - field2
      properties:
        field2:
          type: string
    SchemaС:
      type: object
      required:
        - field3
      properties:
        field3:
          type: string
`)

func Test(t *testing.T) {
	tests.WriteOriginal(t, original)

	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)

	defer tests.CleanGenerated(t, tests.Dir())

	oapi, _ := tests.ReadOAPI(t)
	require.Equal(t, expected, oapi)
}
