package path_replaced

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path GET /hello/{id}
func Hello() {}

var original = tests.TrimOApi(`
openapi: 3.0.3
info:
  title: << YOUR TEXT >>
  version: 0.0.1
paths:
  /hello/{param}:
components:
  schemas:
`)

func Test(t *testing.T) {
	tests.WriteOriginal(t, original)

	_, oapi := tests.ReadOAPI(t)
	_, ok := oapi["paths"].(map[string]any)["/hello/{param}"]
	require.True(t, ok)

	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	_, oapi = tests.ReadOAPI(t)
	_, ok = oapi["paths"].(map[string]any)["/hello/{id}"]
	require.True(t, ok)
}
