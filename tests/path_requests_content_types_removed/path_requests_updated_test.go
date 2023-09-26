package path_replaced

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path POST /hello
//go:generate oapi request text/plain Input
//go:generate oapi request text/plain Input2
//go:generate oapi request Input3
func Hello() {}

var original = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: my oapi
  version: 0.0.1
paths:
  /hello:
    post:
      summary: hello
      requestBody:
        description: request
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Input'
          application/octet-stream:
            schema:
              $ref: '#/components/schemas/Input100'
      responses:
components:
  schemas:
`)

var expected = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: my oapi
  version: 0.0.1
paths:
  /hello:
    post:
      summary: hello
      requestBody:
        description: request
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/Input3'
          text/plain:
            schema:
              oneOf:
                - $ref: '#/components/schemas/Input'
                - $ref: '#/components/schemas/Input2'
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
