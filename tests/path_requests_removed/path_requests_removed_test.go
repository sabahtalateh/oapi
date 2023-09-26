package path_replaced

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path POST /hello
//go:generate oapi request Input
func Hello() {}

//go:generate oapi path POST /goodbye
func goodbye() {}

var original = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
  /goodbye:
    post:
      summary: << YOUR TEXT >>
      requestBody:
        description: << YOUR TEXT >>
        required: true
        content:
          application/json:
            schema:
              oneOf:
                - $ref: '#/components/schemas/Input'
                - $ref: '#/components/schemas/Input2'
      responses:
  /hello:
    post:
      summary: hello url
      requestBody:
        description: hello request
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Input'
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
  /goodbye:
    post:
      summary: << YOUR TEXT >>
      responses:
  /hello:
    post:
      summary: hello url
      requestBody:
        description: hello request
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Input'
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
