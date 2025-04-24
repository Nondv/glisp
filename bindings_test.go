package main

import (
	"github.com/stretchr/testify/require"
	"testing"

	. "nondv.io/glisp/types"
)

func TestBase(t *testing.T) {
	bindings := BuildBaseBindings()

	val, err := ReadEval(bindings, "nil")
	require.NoError(t, err)
	require.True(t, val.IsEmptyList())

	// val, err = ReadEval(bindings, "t")
	// require.NoError(t, err)
	// require.True(t, val.IsSymbol())
	// require.Equal(t, "t", *val.Value.(*string))

	val, err = ReadEval(bindings, "a")
	require.NotNil(t, err)

	bindings = bindings.Assoc(BuildSymbol("a"), BuildInteger(123))
	val, err = ReadEval(bindings, "a")
	require.NoError(t, err)
	require.True(t, val.IsInteger())
	require.Equal(t, 123, val.Value.(int))
}
