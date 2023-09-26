package path_added

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema A
type A struct {
	Field1 string     `json:"field1"`
	Field2 *time.Time `json:"field2"`
}

var expected = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
components:
  schemas:
    A:
      type: object
      required:
        - field1
      properties:
        field1:
          type: string
        field2:
          type: string
          format: date-time
`)

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	oapi, _ := tests.ReadOAPI(t)
	require.Equal(t, expected, oapi)
}
