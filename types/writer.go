package types

import (
	"fmt"
	"strconv"
)

func (v *Value) PrintStr() string {
	if v.IsSymbol() {
		return v.SymbolName()
	}

	if v.IsInteger() {
		return strconv.Itoa(v.ToInt())
	}

	if v.IsEmptyList() {
		return "()"
	}

	if v.IsList() {
		res := "("
		res += v.Car().PrintStr()
		iter := v.Cdr()
		for !iter.IsEmptyList() {
			res += " " + iter.Car().PrintStr()
			iter = iter.Cdr()
		}
		res += ")"
		return res
	}

	if v.IsCons() {
		return fmt.Sprintf("(%v . %s)", v.Car().PrintStr(), v.Cdr().PrintStr())
	}

	if v.IsNativeFn() {
		return "<native fn>"
	}

	panic("Can't convert to string")
}
