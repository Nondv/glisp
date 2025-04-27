package interpreter

import (
	"container/list"
	"errors"
	"os"

	. "nondv.io/glisp/types"
)

/*
 * These are core language functions
 */

func nativeEqual(bindings *Bindings, args *Value) (*Value, error) {
	if args.ListLength() != 2 {
		return nil, errors.New("= requires 2 arguments")
	}
	args, err := evalArgs(bindings, args)
	if err != nil {
		return nil, err
	}

	a := args.Car()
	b := args.Cdr().Car()

	if Equal(a, b) {
		return BuildSymbol("t"), nil
	} else {
		return BuildEmptyList(), nil
	}
}

func nativeCar(bindings *Bindings, args *Value) (*Value, error) {
	args, err := evalArgs(bindings, args)
	if err != nil {
		return nil, err
	}

	argument, err := requireOneArg(args)
	if err != nil {
		return nil, err
	}

	if !argument.IsCons() && !argument.IsEmptyList() {
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

	if !argument.IsCons() && !argument.IsEmptyList() {
		return nil, errors.New("Not a cons cell")
	}

	return argument.Cdr(), nil
}

func nativeCons(bindings *Bindings, args *Value) (*Value, error) {
	if args.ListLength() != 2 {
		return nil, errors.New("cons requires 2 arguments")
	}

	car, err := Eval(bindings, args.Car())
	if err != nil {
		return nil, err
	}

	cdr, err := Eval(bindings, args.Cdr().Car())
	if err != nil {
		return nil, err
	}

	return BuildCons(car, cdr), nil
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

func nativeIf(bindings *Bindings, args *Value) (*Value, error) {
	if args.ListLength() != 3 {
		return nil, errors.New("if requires 3 arguments: condition, then, else")
	}

	condition := args.Car()
	thenBranch := args.Cdr().Car()
	elseBranch := args.Cdr().Cdr().Car()

	conditionVal, err := Eval(bindings, condition)
	if err != nil {
		return nil, err
	}

	if !conditionVal.IsEmptyList() {
		return Eval(bindings, thenBranch)
	}

	return Eval(bindings, elseBranch)
}

func nativePrint(bindings *Bindings, args *Value) (*Value, error) {
	args, err := evalArgs(bindings, args)
	if err != nil {
		return nil, err
	}

	lastValue := BuildEmptyList()
	for iter := args; !iter.IsEmptyList(); iter = args.Cdr() {
		lastValue = iter.Car()
		Print(lastValue)
	}

	return lastValue, nil
}

func nativeDefine(bindings *Bindings, args *Value) (*Value, error) {
	if args.ListLength() != 2 {
		return nil, errors.New("define requires 2 arguments")
	}

	sym := args.Car()
	if !sym.IsSymbol() {
		return nil, errors.New("syntax: (define SYMBOL SEXP) ")
	}

	sexp := args.Cdr().Car()

	value, err := Eval(bindings, sexp)
	if err != nil {
		return nil, err
	}

	// can't just assign directly because the head would be pointing at itself
	// so first create a copy so the new head points at that
	bindingsCopy := *bindings
	*bindings = *bindingsCopy.Assoc(sym, value)

	return value, nil
}

func nativeLoad(bindings *Bindings, args *Value) (*Value, error) {
	if args.ListLength() != 1 {
		return nil, errors.New("load requires 1 argument")
	}
	argument, err := Eval(bindings, args.Car())
	if err != nil {
		return nil, err
	}
	if !argument.IsString() {
		return nil, errors.New("load requires a string as its argument")
	}

	contents, err := os.ReadFile(argument.ToStr())
	if err != nil {
		return nil, err
	}

	return ReadEvalAll(bindings, string(contents))
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
