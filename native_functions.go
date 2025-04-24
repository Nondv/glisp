package main

import (
	"container/list"
	"errors"

	. "nondv.io/glisp/types"
)

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
	argument, err := requireOneArg(args)
	if err != nil {
		return nil, err
	}

	result, err := Eval(bindings, argument)

	return result, err

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
