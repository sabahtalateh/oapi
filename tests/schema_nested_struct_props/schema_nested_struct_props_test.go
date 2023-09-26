package path_added

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/sabahtalateh/oapi/tests"
)

// TODO

//go:generate oapi schema A
type A struct {
	Nested struct {
		Field string `json:"field"`
	} `json:"nested"`
}

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.Contains(t, err.Error(), "nested struct not supported: nested")

	defer tests.CleanGenerated(t, tests.Dir())
}
