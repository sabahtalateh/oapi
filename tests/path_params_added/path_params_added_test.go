package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path GET /hello/{id}/{param}
func GetHello() {}

var expected = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
  /hello/{id}/{param}:
    get:
      summary: << YOUR TEXT >>
      parameters:
        - name: id
          in: path
          description: << YOUR TEXT >>
          required: true
          schema:
            type: string
        - name: param
          in: path
          description: << YOUR TEXT >>
          required: true
          schema:
            type: string
      responses:
components:
  schemas:
`)

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	oapi, _ := tests.ReadOAPI(t)
	require.Equal(t, expected, oapi)
}
