package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

// TODO finish test

//go:generate oapi path GET /hello/{p1}/{p2}/{p3}
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
          description: Descr 1
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
        - name: param1
          in: path
          description: Descr 2
          required: true
          schema:
            type: string
        - name: param2
          in: path
          description: Descr 3
          required: true
          schema:
            type: string
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
  /hello/{p1}/{p2}/{p3}:
    get:
      summary: << YOUR TEXT >>
      parameters:
        - name: p1
          in: path
          description: << YOUR TEXT >>
          required: true
          schema:
            type: string
        - name: p2
          in: path
          description: << YOUR TEXT >>
          required: true
          schema:
            type: string
        - name: p3
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
