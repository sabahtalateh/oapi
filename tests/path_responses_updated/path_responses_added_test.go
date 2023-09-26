package path_replaced

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path POST /hello
//go:generate oapi response 200 Output
//go:generate oapi response 200 text/plain Output
//go:generate oapi response 200 text/plain Output2
//go:generate oapi response default Error
//go:generate oapi response default text/plain Error
func Hello() {}

var expected = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
  /hello:
    post:
      summary: << YOUR TEXT >>
      responses:
        200:
          description: << YOUR TEXT >>
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/Output'
            text/plain:
              schema:
                oneOf:
                  - $ref: '#/components/schemas/Output'
                  - $ref: '#/components/schemas/Output2'
        default:
          description: << YOUR TEXT >>
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
	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	oapi, _ := tests.ReadOAPI(t)
	require.Equal(t, expected, oapi)
}
