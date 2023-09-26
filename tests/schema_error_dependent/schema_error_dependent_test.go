package path_added

import (
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/sabahtalateh/oapi/tests"
)

//go:generate oapi schema A
type A struct {
	B
}

type B struct {
}

func Test(t *testing.T) {
	err := tests.GoGenerate(t, tests.Dir())
	require.Contains(t, err.Error(), "can not be used as part of type A")
	require.Contains(t, err.Error(), "not marked with `//go:generate oapi schema ..`")
	defer tests.CleanGenerated(t, tests.Dir())
}
