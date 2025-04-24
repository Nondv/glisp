package main

import (
	"bufio"
	"errors"
	"os"

	"nondv.io/glisp/reader"
	. "nondv.io/glisp/types"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		print("> ")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()
		result, err := ReadEval(BuildBaseBindings(), input)
		if err != nil {
			println("Err!")
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
		if !parameter.IsSymbol() {
			return nil, errors.New("format: (lambda SYMBOL BODY)")
		}

		lambdaBindings := bindings.Assoc(parameter, args)
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

func BuildBaseBindings() *Bindings {
	result := &Bindings{"nil", BuildEmptyList(), nil}
	// result = result.Assoc(BuildSymbol("t"), BuildSymbol("t"))
	result = result.Assoc(BuildSymbol("eval"), BuildNativeFn(nativeEval))
	result = result.Assoc(BuildSymbol("+"), BuildNativeFn(nativePlus))
	result = result.Assoc(BuildSymbol("car"), BuildNativeFn(nativeCar))
	result = result.Assoc(BuildSymbol("cdr"), BuildNativeFn(nativeCdr))

	return result
}
