package path_replaced

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path POST /hello
//go:generate oapi response 200 Output
//go:generate oapi response 201 text/plain Output
//go:generate oapi response 202 text/plain Output2
//go:generate oapi response 202 text/plain Output3
//go:generate oapi response default Error
//go:generate oapi response default text/plain Error
func Hello() {}

var original = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
  /hello:
    post:
      summary: fiz
      responses:
        201:
          description: 201
          content:
            text/plain:
              schema:
                $ref: '#/components/schemas/Output'
        default:
          description: default
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
            text/plain:
              schema:
                $ref: '#/components/schemas/Error'
        
        202:
          description: 202
          content:
            text/plain:
              schema:
                oneOf:
                  - $ref: '#/components/schemas/Output2'
                  - $ref: '#/components/schemas/Output3'
        200:
          description: 200
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Output'
        
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
    post:
      summary: fiz
      responses:
        200:
          description: 200
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Output'
        201:
          description: 201
          content:
            text/plain:
              schema:
                $ref: '#/components/schemas/Output'
        202:
          description: 202
          content:
            text/plain:
              schema:
                oneOf:
                  - $ref: '#/components/schemas/Output2'
                  - $ref: '#/components/schemas/Output3'
        default:
          description: default
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Error'
            text/plain:
              schema:
                $ref: '#/components/schemas/Error'
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
