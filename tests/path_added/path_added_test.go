package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path GET /hello
func Hello() {}

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	_, oapi := tests.ReadOAPI(t)
	_, ok := oapi["paths"].(map[string]any)["/hello"]
	require.True(t, ok)
}
