package path_added

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi path GET /hello
type A struct{}

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.Contains(t, err.Error(), "malformed path")
	tests.CleanGenerated(t, tests.Dir())
}
