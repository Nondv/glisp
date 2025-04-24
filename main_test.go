package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	. "nondv.io/glisp/types"
)


func TestLisp(t *testing.T) {
	bindings := BuildBaseBindings()

	require.Equal(t, "()", readEvalPrintNoErr(bindings, "nil"))
	require.Equal(t, "(lambda X 123)", readEvalPrintNoErr(bindings, "(lambda X 123)"))

	require.Equal(t, "<native fn>", readEvalPrintNoErr(bindings, "+"))
	require.Equal(t, "0", readEvalPrintNoErr(bindings, "(+)"))
	require.Equal(t, "6", readEvalPrintNoErr(bindings, "(+ 1 2 3)"))

	require.Equal(t, "bla",
		readEvalPrintNoErr(bindings, "((lambda X (car X)) bla)"))
	require.Equal(t, "123",
		readEvalPrintNoErr(bindings, "((lambda X 123) 456 789)"))
	require.Equal(t, "2",
		readEvalPrintNoErr(bindings, "((lambda X (car (cdr X))) 1 2 3)"))
	require.Equal(t, "579",
		readEvalPrintNoErr(bindings, "((lambda X (+ (car (cdr X)) (car X))) 123 456)"))


	code := "((lambda quote (quote a b c d)) lambda X X)"
	require.Equal(t, "(a b c d)", readEvalPrintNoErr(bindings, code))
}

func TestNoLet(t *testing.T) {
	bindings := BuildBaseBindings()
	withQuote := func(code string) string {
		return fmt.Sprintf("((lambda quote %s) lambda X X)", code)
	}

	code := "(quote a b c d)"
	require.Equal(t, "(a b c d)", readEvalPrintNoErr(bindings, withQuote(code)))
}


func readEvalPrintNoErr(bindings *Bindings, txt string) string {
	res, err := ReadEval(bindings, txt)
	if err != nil {
		println(err.Error())
		panic("err!")
	}
	return res.PrintStr()
}
