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
    post:
components:
  schemas:
`)

// Test when existing methods not removed even it's not in `//go:generate oapi path ..`
func Test(t *testing.T) {
	tests.WriteOriginal(t, original)

	_, api := tests.ReadOAPI(t)
	_, ok := api["paths"].(map[string]any)["/hello"].(map[string]any)["post"]
	require.True(t, ok)

	require.NoError(t, tests.GoGenerate(t, tests.Dir()))
	defer tests.CleanGenerated(t, tests.Dir())

	_, api = tests.ReadOAPI(t)
	_, ok = api["paths"].(map[string]any)["/hello"].(map[string]any)["post"]
	require.True(t, ok)
}
