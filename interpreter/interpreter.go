package interpreter

import (
	"bufio"
	"errors"
	"os"

	"nondv.io/glisp/reader"
	. "nondv.io/glisp/types"
)

func BuildBaseBindings() *Bindings {
	result := &Bindings{"nil", BuildEmptyList(), nil}
	// result = result.Assoc(BuildSymbol("t"), BuildSymbol("t"))
	result = result.Assoc(BuildSymbol("eval"), BuildNativeFn(nativeEval))
	result = result.Assoc(BuildSymbol("let"), BuildNativeFn(nativeLet))
	result = result.Assoc(BuildSymbol("if"), BuildNativeFn(nativeIf))
	result = result.Assoc(BuildSymbol("+"), BuildNativeFn(nativePlus))
	result = result.Assoc(BuildSymbol("car"), BuildNativeFn(nativeCar))
	result = result.Assoc(BuildSymbol("cdr"), BuildNativeFn(nativeCdr))
	result = result.Assoc(BuildSymbol("cons"), BuildNativeFn(nativeCons))
	result = result.Assoc(BuildSymbol("print"), BuildNativeFn(nativePrint))

	return result
}

func Repl(baseBindings *Bindings) {
	println("NOTE: this repl can only read-eval one line at a time")
	println("If you want to eval something more complex, put it in a file and run")
	println("    glisp yourcode.lisp")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		result, err := ReadEval(baseBindings, input)
		if err != nil {
			if _, ok := err.(*reader.NoNextSexpError); !ok {
				println("Err: ", err.Error())
			}
		} else {
			Print(result)
		}
	}
}

func ReadEval(bindings *Bindings, txt string) (*Value, error) {
	sexp, err := reader.Read(txt)
	if err != nil {
		return nil, err
	}

	return Eval(bindings, sexp)
}

// Evals all sexps and returns last value
func ReadEvalAll(bindings *Bindings, txt string) (*Value, error) {
	sexps, err := reader.ReadAll(txt)
	if err != nil {
		return sexps, err
	}

	var lastResult *Value
	for iter := sexps; !iter.IsEmptyList(); iter = iter.Cdr() {
		lastResult, err = Eval(bindings, iter.Car())
		if err != nil {
			return lastResult, err
		}
	}

	return lastResult, nil
}

func Eval(bindings *Bindings, v *Value) (*Value, error) {
	if v.IsInteger() || v.IsEmptyList() {
		return v, nil
	}

	if v.IsSymbol() {
		val, found := bindings.Lookup(v)
		if !found {
			return nil, errors.New("Undefined")
		}
		return val, nil
	}

	if v.IsNativeFn() {
		return nil, errors.New("Not eval-able")
	}

	if v.IsList() {
		fn := v.Car()
		if fn.IsSymbol() && fn.SymbolName() == "lambda" {
			return v, nil
		}

		return callFn(bindings, fn, v.Cdr())
	}

	panic("Unexpected eval argument")
}

func Print(v *Value) {
	println(v.PrintStr())
}

func callFn(bindings *Bindings, fn *Value, args *Value) (*Value, error) {
	fn, err := Eval(bindings, fn)
	if err != nil {
		return nil, err
	}

	if fn.IsNativeFn() {
		resPointer, err := fn.NativeFn()(bindings, args)
		return resPointer, err
	}

	if fn.IsList() && fn.Car().IsLambdaSymbol() {
		parameter := fn.Cdr().Car()
		if !parameter.IsSymbol() && !parameter.IsList() {
			return nil, errors.New("format: (lambda SYMBOL-OR-LIST BODY)")
		}
		var lambdaBindings *Bindings
		if parameter.IsSymbol() {
			lambdaBindings = bindings.Assoc(parameter, args)
		} else {
			if parameter.ListLength() != args.ListLength() {
				return nil, errors.New("too many/not enough arguments")
			}

			lambdaBindings = bindings
			for iter := parameter; !iter.IsEmptyList(); {
				varSym := iter.Car()
				if !varSym.IsSymbol() {
					return nil, errors.New("parameter is not a symbol")
				}

				argN, err := Eval(bindings, args.Car())
				if err != nil {
					return nil, err
				}

				lambdaBindings = lambdaBindings.Assoc(varSym, argN)
				iter = iter.Cdr()
				args = args.Cdr()
			}
		}
		res := BuildEmptyList()
		body := fn.Cdr().Cdr()
		for !body.IsEmptyList() {
			res, err = Eval(lambdaBindings, body.Car())
			if err != nil {
				return nil, err
			}
			body = body.Cdr()
		}
		return res, nil
	}

	return nil, errors.New("Not a function")
}
