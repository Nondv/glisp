package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"nondv.io/glisp/reader"
	. "nondv.io/glisp/types"
	. "nondv.io/glisp/types/bindings"
)

func main() {
	inputBytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic("!!!!")
	}

	fmt.Println(ReadEval(BuildBaseBindings(), string(inputBytes)))
}

func ReadEval(bindings *Bindings, txt string) (Value, error) {
	sexp, err := reader.Read(txt)
	if err != nil {
		return sexp, err
	}

	return Eval(bindings, sexp)
}


func Eval(bindings *Bindings, v Value) (Value, error) {
	if IsInteger(v) || IsEmptyList(v) {
		return v, nil
	}

	if IsSymbol(v) {
		val, found := bindings.Lookup(v)
		if !found {
			return Value{}, errors.New("Undefined")
		}
		return val, nil
	}

	if isList(v) {
		fn := (*v.Value.(*Cons)).Car
		args := (*v.Value.(*Cons)).Cdr
		return callFn(bindings, fn, args)
	}

	panic("Unexpected eval argument")
}


func callFn(bindings *Bindings, fn Value, args Value) (Value, error) {
	if isList(fn) && isLambdaSym(listFirst(fn)) {
		panic("TODO: call the lambda")
	}

	fn, err := Eval(bindings, fn)
	if err != nil {
		return Value{}, err
	}

	return args, nil
}

func isList(v Value) bool {
	iter := v
	for IsCons(iter) {
		iter = (*iter.Value.(*Cons)).Cdr
	}

	return IsEmptyList(v)
}

func listFirst(list Value) Value {
	return (*list.Value.(*Cons)).Car
}

func isLambdaSym(v Value) bool {
	return IsSymbol(v) && (*v.Value.(*string)) == "lambda"
}


// func Equals(a Value, b Value) bool {
// 	a = Eval(a)
// 	b = Eval(b)

// 	if a.ValueType != b.ValueType {
// 		return false
// 	}

// 	if IsEmptyList(a) {
// 		return true
// 	}

// 	if IsSymbol(a) {
// 		return *(a.Value.(*string)) == *(b.Value.(*string))
// 	}

// 	if IsInteger(a) {
// 		return a.Value.(int) == b.Value.(int)
// 	}

// 	if IsList

// 	panic("Unexpected args type for Equals")
// }
