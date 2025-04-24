package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"nondv.io/glisp/interpreter"
	. "nondv.io/glisp/types"
)


func TestLisp(t *testing.T) {
	bindings := interpreter.BuildBaseBindings()

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


	code = `(let ((double (lambda X (+ (car X) (car X)))))
                  (double 123))`
	require.Equal(t, "246", readEvalPrintNoErr(bindings, code))

	code = `((lambda X (eval (car X))) (+ 123 111))`
	require.Equal(t, "234", readEvalPrintNoErr(bindings, code))

	code = `(let ((double (lambda X
                                (let ((x (eval (car X))))
                                   (+ x x))))
                      (quadriple (lambda X
                                   (let ((x (eval (car X))))
                                     (double (double x))))))
                  (quadriple (+ 1 2 3)))`
	require.Equal(t, "24", readEvalPrintNoErr(bindings, code))

	// if parameters are in a list, arguments are evaluated
	// automatically before entering the function
	code = "((lambda (x) (+ x x)) (+ 1 2 3))"
	require.Equal(t, "12", readEvalPrintNoErr(bindings, code))
}

func TestNoLet(t *testing.T) {
	bindings := interpreter.BuildBaseBindings()
	withQuote := func(code string) string {
		return fmt.Sprintf("((lambda quote %s) lambda X X)", code)
	}

	code := "(quote a b c d)"
	require.Equal(t, "(a b c d)", readEvalPrintNoErr(bindings, withQuote(code)))
}


func readEvalPrintNoErr(bindings *Bindings, txt string) string {
	res, err := interpreter.ReadEval(bindings, txt)
	if err != nil {
		println(err.Error())
		panic("err!")
	}
	return res.PrintStr()
}
