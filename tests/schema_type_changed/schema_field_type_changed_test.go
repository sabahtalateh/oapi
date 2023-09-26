package path_added

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema SchemaA
type A []B

//go:generate oapi schema SchemaB
type B struct {
	Prop string `json:"prop"`
}

//go:generate oapi schema SchemaA2
type A2 struct {
	Field1 time.Time `json:"field1"`
}

//go:generate oapi schema SchemaB2
type B2 struct {
	Prop *string `json:"prop"`
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
          type: string
      description: hehe1
    SchemaA2:
      type: array
      items:
        $ref: '#/components/schemas/SchemaB2'
      description: hehe2
    SchemaB:
      type: object
      properties:
        prop:
          type: string
    SchemaB2:
      type: object
      required:
        - prop
      properties:
        prop:
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
      type: array
      items:
        $ref: '#/components/schemas/SchemaB'
      description: hehe1
    SchemaA2:
      type: object
      required:
        - field1
      properties:
        field1:
          type: string
          format: date-time
      description: hehe2
    SchemaB:
      type: object
      required:
        - prop
      properties:
        prop:
          type: string
    SchemaB2:
      type: object
      properties:
        prop:
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
