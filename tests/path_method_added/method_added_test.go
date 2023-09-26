package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path GET /hello
func GetHello() {}

//go:generate oapi path HEAD /hello
func HeadHello() {}

//go:generate oapi path POST /hello
func PostHello() {}

//go:generate oapi path PUT /hello
func PutHello() {}

//go:generate oapi path DELETE /hello
func DeleteHello() {}

//go:generate oapi path CONNECT /hello
func ConnectHello() {}

//go:generate oapi path OPTIONS /hello
func OptionsHello() {}

//go:generate oapi path TRACE /hello
func TraceHello() {}

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.NoError(t, err)
	defer tests.CleanGenerated(t, tests.Dir())

	_, oapi := tests.ReadOAPI(t)
	_, ok := oapi["paths"].(map[string]any)["/hello"].(map[string]any)["get"]
	require.True(t, ok)

	_, ok = oapi["paths"].(map[string]any)["/hello"].(map[string]any)["head"]
	require.True(t, ok)

	_, ok = oapi["paths"].(map[string]any)["/hello"].(map[string]any)["post"]
	require.True(t, ok)

	_, ok = oapi["paths"].(map[string]any)["/hello"].(map[string]any)["put"]
	require.True(t, ok)

	_, ok = oapi["paths"].(map[string]any)["/hello"].(map[string]any)["delete"]
	require.True(t, ok)

	_, ok = oapi["paths"].(map[string]any)["/hello"].(map[string]any)["connect"]
	require.True(t, ok)

	_, ok = oapi["paths"].(map[string]any)["/hello"].(map[string]any)["options"]
	require.True(t, ok)

	_, ok = oapi["paths"].(map[string]any)["/hello"].(map[string]any)["trace"]
	require.True(t, ok)

	require.True(t, ok)
}
