package main

import (
	"errors"

	"nondv.io/glisp/interpreter"
	. "nondv.io/glisp/types"
)

func main() {
	bindings := interpreter.BuildBaseBindings()
	bindings = bindings.Assoc(BuildSymbol("sqr"), BuildNativeFn(nativeSqr))

	interpreter.Repl(bindings)
}

// example of extending the language
func nativeSqr(bindings *Bindings, args *Value) (*Value, error) {
	if args.ListLength() != 1 {
		return nil, errors.New("Only one argument expected")
	}

	if !args.Car().IsInteger() {
		return nil, errors.New("Integer expected")
	}

	n := args.Car().ToInt()
	return BuildInteger(n * n), nil
}
