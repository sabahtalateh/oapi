package path_added

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema SchemaA
type A struct {
	Field1  []string `json:"field1"`
	Field2  B        `json:"field2"`
	Field3  []B      `json:"field3"`
	Field4  string   `json:"field4"`
	Field5  B        `json:"field5"`
	Field6  []B      `json:"field6"`
	Field7  string   `json:"field7"`
	Field8  []string `json:"field8"`
	Field9  []B      `json:"field9"`
	Field10 string   `json:"field10"`
	Field11 []string `json:"field11"`
	Field12 B        `json:"field12"`
}

//go:generate oapi schema SchemaB
type B struct {
	Field2 string `json:"field2"`
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
      type: object
      properties:
        field1:
          type: number
          format: double
          description: Something1
        field2:
          type: number
          format: double
          description: Something2
        field3:
          type: number
          format: double
          description: Something3
        field4:
          type: array
          items:
            type: string
          description: Something4
        field5:
          type: array
          items:
            type: string
          description: Something5
        field6:
          type: array
          items:
            type: string
          description: Something6
        field7:
          $ref: '#/components/schemas/SchemaB'
          description: Something7
        field8:
          $ref: '#/components/schemas/SchemaB'
          description: Something8
        field9:
          $ref: '#/components/schemas/SchemaB'
          description: Something9
        field10:
          type: array
          items:
            $ref: '#/components/schemas/SchemaB'
          description: Something10
        field11:
          type: array
          items:
            $ref: '#/components/schemas/SchemaB'
          description: Something11
        field12:
          type: array
          items:
            $ref: '#/components/schemas/SchemaB'
          description: Something12
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
      type: object
      required:
        - field1
        - field2
        - field3
        - field4
        - field5
        - field6
        - field7
        - field8
        - field9
        - field10
        - field11
        - field12
      properties:
        field1:
          type: array
          items:
            type: string
          description: Something1
        field2:
          $ref: '#/components/schemas/SchemaB'
          description: Something2
        field3:
          type: array
          items:
            $ref: '#/components/schemas/SchemaB'
          description: Something3
        field4:
          type: string
          description: Something4
        field5:
          $ref: '#/components/schemas/SchemaB'
          description: Something5
        field6:
          type: array
          items:
            $ref: '#/components/schemas/SchemaB'
          description: Something6
        field7:
          type: string
          description: Something7
        field8:
          type: array
          items:
            type: string
          description: Something8
        field9:
          type: array
          items:
            $ref: '#/components/schemas/SchemaB'
          description: Something9
        field10:
          type: string
          description: Something10
        field11:
          type: array
          items:
            type: string
          description: Something11
        field12:
          $ref: '#/components/schemas/SchemaB'
          description: Something12
    SchemaB:
      type: object
      required:
        - field2
      properties:
        field2:
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
