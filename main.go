package main

import (
	"errors"
	"os"

	"nondv.io/glisp/interpreter"
	. "nondv.io/glisp/types"
)

func main() {
	bindings := interpreter.BuildBaseBindings()
	bindings = bindings.Assoc(BuildSymbol("sqr"), BuildNativeFn(nativeSqr))

	// No arguments provided
	if len(os.Args) == 1 {
		interpreter.Repl(bindings)
		return
	}

	filename := os.Args[1]
	contents, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	lastResult, err := interpreter.ReadEvalAll(bindings, string(contents))
	if err != nil {
		panic(err)
	}

	interpreter.Print(lastResult)
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
