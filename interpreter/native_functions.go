package interpreter

import (
	"container/list"
	"errors"

	. "nondv.io/glisp/types"
)

/*
 * These are core language functions
 */

func nativeCar(bindings *Bindings, args *Value) (*Value, error) {
	args, err := evalArgs(bindings, args)
	if err != nil {
		return nil, err
	}

	argument, err := requireOneArg(args)
	if err != nil {
		return nil, err
	}

	if !argument.IsCons() {
		return nil, errors.New("Not a cons cell")
	}

	return argument.Car(), nil
}

func nativeCdr(bindings *Bindings, args *Value) (*Value, error) {
	args, err := evalArgs(bindings, args)
	if err != nil {
		return nil, err
	}

	argument, err := requireOneArg(args)
	if err != nil {
		return nil, err
	}

	if !argument.IsCons() {
		return nil, errors.New("Not a cons cell")
	}

	return argument.Cdr(), nil
}

func nativePlus(bindings *Bindings, args *Value) (*Value, error) {
	args, err := evalArgs(bindings, args)
	if err != nil {
		return nil, err
	}

	res := 0
	for iter := args; iter.IsCons(); iter = iter.Cdr() {
		arg := iter.Car()
		if !arg.IsInteger() {
			return nil, errors.New("non-integer argument")
		}

		res += arg.ToInt()
	}

	return BuildInteger(res), nil
}

func nativeEval(bindings *Bindings, args *Value) (*Value, error) {
	args, err := evalArgs(bindings, args)
	if err != nil {
		return nil, err
	}

	argument, err := requireOneArg(args)
	if err != nil {
		return nil, err
	}

	return Eval(bindings, argument)

}

func nativeLet(bindings *Bindings, args *Value) (*Value, error) {
	varList := args.Car()
	body := args.Cdr()
	if !varList.IsList() {
		return nil, errors.New("Expected a varlist")
	}

	newBindings := bindings
	for iter := varList; !iter.IsEmptyList(); iter = iter.Cdr() {
		declaration := iter.Car()
		if !declaration.IsList() && declaration.ListLength() != 2 {
			return nil, errors.New("invalid varlist")
		}
		varSym := declaration.Car()
		if !varSym.IsSymbol() {
			return nil, errors.New("vars must be symbols")
		}
		value, err := Eval(newBindings, declaration.Cdr().Car())
		if err != nil {
			return nil, err
		}
		newBindings = newBindings.Assoc(varSym, value)
	}

	result := BuildEmptyList()
	for iter := body; !iter.IsEmptyList(); iter = iter.Cdr() {
		var err error
		result, err = Eval(newBindings, iter.Car())
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func evalArgs(bindings *Bindings, args *Value) (*Value, error) {
	if !args.IsList() {
		panic("args aren't a list for some reason")
	}

	values := list.New()

	for iter := args; !iter.IsEmptyList(); iter = iter.Cdr() {
		v, err := Eval(bindings, iter.Car())
		if err != nil {
			return nil, err
		}
		values.PushFront(v)
	}

	result := BuildEmptyList()
	for e := values.Front(); e != nil; e = e.Next() {
		result = BuildCons(e.Value.(*Value), result)
	}

	return result, nil
}

func requireOneArg(args *Value) (*Value, error) {
	if !args.IsList() {
		panic("args aren't a list for some reason")
	}

	if args.ListLength() != 1 {
		return nil, errors.New("Only one argument expected")
	}

	return args.Car(), nil
}
