package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path GET /hello/{param1}/{param2}/{param3}
func GetHello() {}

var original = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
  /hello/{param1}/{param2}/{param3}:
    get:
      summary: << YOUR TEXT >>
      parameters:
        - name: param3
          in: path
          description: << YOUR TEXT >>
          required: true
          schema:
            type: string
        - name: param1
          in: path
          description: << YOUR TEXT >>
          required: true
          schema:
            type: string
        - name: param2
          in: path
          description: << YOUR TEXT >>
          required: true
          schema:
            type: string
        - name: query_param
          in: query
          description: << YOUR TEXT >>
          required: false
          schema:
            type: string
            minLength: 999999999
      responses:
components: 
  schemas:
`)

var expected = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
  /hello/{param1}/{param2}/{param3}:
    get:
      summary: << YOUR TEXT >>
      parameters:
        - name: param1
          in: path
          description: << YOUR TEXT >>
          required: true
          schema:
            type: string
        - name: param2
          in: path
          description: << YOUR TEXT >>
          required: true
          schema:
            type: string
        - name: param3
          in: path
          description: << YOUR TEXT >>
          required: true
          schema:
            type: string
        - name: query_param
          in: query
          description: << YOUR TEXT >>
          required: false
          schema:
            type: string
            minLength: 999999999
      responses:
components: 
  schemas:
`)

func Test(t *testing.T) {
	tests.WriteOriginal(t, original)
	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	oapi, _ := tests.ReadOAPI(t)
	require.Equal(t, expected, oapi)
}
