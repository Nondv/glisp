package reader

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSymbol(t *testing.T) {
	requireSymbol(t, "hello", Read("hello"))
	requireSymbol(t, "hello", Read(" hello   "))
	requireSymbol(t, "abc+def-ghi/123", Read("abc+def-ghi/123"))
}

func TestInteger(t *testing.T) {
	requireInteger(t, 123, Read("123"))
	requireInteger(t, 456, Read("456 "))
	requireInteger(t, -9, Read("-9"))
	requireInteger(t, 0, Read("000"))
	requireInteger(t, 0, Read("-0"))

	requireSymbol(t, "--123", Read(" --123"))
}

func requireSymbol(t *testing.T, name string, val Value) {
	require.Equal(t, SymbolReference, val.valueType)
	strPointer, ok := val.value.(*string)

	require.True(t, ok)
	require.Equal(t, name, *strPointer)
}

func requireInteger(t *testing.T, expected int, val Value) {
	require.Equal(t, IntegerReference, val.valueType)
	actual, ok := val.value.(int)

	require.True(t, ok)
	require.Equal(t, expected, actual)
}
