package path_added

import (
	"github.com/sabahtalateh/oapi/tests"
	"github.com/stretchr/testify/require"
	"testing"
)

//go:generate oapi path TAKE /hello
func Hello() {}

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	str := err.Error()
	require.Contains(t, str, "unknown method: TAKE")

	defer tests.CleanGenerated(t, tests.Dir())
}
