package path_added

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema A
func A() {}

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.Contains(t, err.Error(), "malformed schema")
	tests.CleanGenerated(t, tests.Dir())
}
