package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

// TODO finish test

//go:generate oapi path GET /hello/{param}
func GetHello() {}

var original = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
  /hello/{param}:
    get:
      summary: << YOUR TEXT >>
      parameters:
        - name: query_param1
          in: query
          description: Descr 1
          required: false
          schema:
            type: string
            minLength: 999999999
        - name: param
          in: path
          description: Descr 2
          required: true
          schema:
            type: string
        - name: query_param2
          in: query
          description: Descr 3
          required: false
          schema:
            type: string
            minLength: -1
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
  /hello/{param}:
    get:
      summary: << YOUR TEXT >>
      parameters:
        - name: param
          in: path
          description: Descr 2
          required: true
          schema:
            type: string
        - name: query_param1
          in: query
          description: Descr 1
          required: false
          schema:
            type: string
            minLength: 999999999
        - name: query_param2
          in: query
          description: Descr 3
          required: false
          schema:
            type: string
            minLength: -1
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
