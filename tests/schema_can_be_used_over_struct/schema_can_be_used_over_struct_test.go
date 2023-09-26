package path_added

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema A
func A() {}

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.Contains(t, err.Error(), "malformed schema")
	tests.CleanGenerated(t, tests.Dir())
}
