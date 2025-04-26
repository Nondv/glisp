package reader

import (
	"github.com/stretchr/testify/require"
	"testing"

	. "nondv.io/glisp/types"
)

func TestSymbol(t *testing.T) {
	requireSymbol(t, "hello", readNoErr("hello"))
	requireSymbol(t, "hello", readNoErr(" hello   "))
	requireSymbol(t, "abc+def-ghi/123", readNoErr("abc+def-ghi/123"))

	// These aren't reserved
	requireSymbol(t, "nil", readNoErr("nil"))
	requireSymbol(t, "t", readNoErr("t"))
}

func TestInteger(t *testing.T) {
	requireInteger(t, 123, readNoErr("123"))
	requireInteger(t, 456, readNoErr("456 "))
	requireInteger(t, -9, readNoErr("-9"))
	requireInteger(t, 0, readNoErr("000"))
	requireInteger(t, 0, readNoErr("-0"))

	requireSymbol(t, "--123", readNoErr(" --123"))
}

func TestList(t *testing.T) {
	requireEmptyList(t, readNoErr("()"))
	requireEmptyList(t, readNoErr("(    \n   )"))

	value := readNoErr("(())")
	requireCons(t, value)
	consPtr, _ := value.Value.(*Cons)
	cons := *consPtr
	requireEmptyList(t, cons.Car)
	requireEmptyList(t, cons.Cdr)

	value = readNoErr("(a   b(c))")
	requireCons(t, value)
	consPtr, _ = value.Value.(*Cons)
	cons = *consPtr
	requireSymbol(t, "a", cons.Car)

	requireCons(t, cons.Cdr)
	consPtr, _ = cons.Cdr.Value.(*Cons)
	cons = *consPtr
	requireSymbol(t, "b", cons.Car)

	requireCons(t, cons.Cdr)
	consPtr, _ = cons.Cdr.Value.(*Cons)
	cons = *consPtr
	// nested list (c)
	requireCons(t, cons.Car)
	consPtr, _ = cons.Car.Value.(*Cons)
	requireSymbol(t, "c", (*consPtr).Car)
	requireEmptyList(t, (*consPtr).Cdr)

	requireEmptyList(t, cons.Cdr)
}

func TestReadAll(t *testing.T) {
	sexps, err := ReadAll("(hello-world)")
	require.NoError(t, err)
	require.Equal(t, 1, sexps.ListLength())

	sexps, err = ReadAll("")
	require.NoError(t, err)
	require.Equal(t, 0, sexps.ListLength())

	sexps, err = ReadAll("(1 (2 3)) (4) 5")
	require.NoError(t, err)
	require.Equal(t, 3, sexps.ListLength())
	requireInteger(t, 5, sexps.Cdr().Cdr().Car())
	requireInteger(t, 1, sexps.Car().Car())
	require.Equal(t, 2, sexps.Car().Cdr().Car().ListLength())

	sexps, err = ReadAll("(()")
	require.NotNil(t, err)

	sexps, err = ReadAll("()(")
	require.NotNil(t, err)
}

func requireEmptyList(t *testing.T, val *Value) {
	require.True(t, val.IsEmptyList())
	require.Nil(t, val.Value)
}

func requireCons(t *testing.T, val *Value) {
	require.True(t, val.IsCons())
	_, ok := val.Value.(*Cons)
	require.True(t, ok)
}

func requireSymbol(t *testing.T, name string, val *Value) {
	require.True(t, val.IsSymbol())
	strPointer, ok := val.Value.(*string)

	require.True(t, ok)
	require.Equal(t, name, *strPointer)
}

func requireInteger(t *testing.T, expected int, val *Value) {
	require.True(t, val.IsInteger())
	actual, ok := val.Value.(int)

	require.True(t, ok)
	require.Equal(t, expected, actual)
}


func readNoErr(txt string) *Value {
	val, _ := Read(txt)
	return val
}
