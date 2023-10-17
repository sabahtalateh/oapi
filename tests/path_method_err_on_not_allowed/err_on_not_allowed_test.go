package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path TAKE /hello
func Hello() {}

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	str := err.Error()
	require.Contains(t, str, "unknown method: TAKE")

	defer tests.CleanGenerated(t, tests.Dir())
}
