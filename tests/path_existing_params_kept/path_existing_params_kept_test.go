package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path GET /hello
func GetHello() {}

var original = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
  /hello:
    get:
      summary: << YOUR TEXT >>
      parameters:
        - name: param
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
  /hello:
    get:
      summary: << YOUR TEXT >>
      parameters:
        - name: param
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

func TestExistingParamsKept(t *testing.T) {
	tests.WriteOriginal(t, original)
	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	oapi, _ := tests.ReadOAPI(t)
	require.Equal(t, expected, oapi)
}
