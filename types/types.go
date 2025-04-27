package types

const (
	symbolReference    = "sym"
	integerReference   = "int"
	consReference      = "cons"
	emptyListReference = "()"
	nativeFnReference  = "<native fn>"
	stringReference    = "string"
)

type Value struct {
	ValueType string
	Value     any
}

type Cons struct {
	Car *Value
	Cdr *Value
}

func BuildSymbol(name string) *Value {
	return &Value{symbolReference, name}
}

func BuildInteger(n int) *Value {
	return &Value{integerReference, n}
}

func BuildCons(car *Value, cdr *Value) *Value {
	return &Value{consReference, &Cons{car, cdr}}
}

func BuildEmptyList() *Value {
	return &Value{emptyListReference, nil}
}

func BuildNativeFn(f func(*Bindings, *Value) (*Value, error)) *Value {
	return &Value{nativeFnReference, f}
}

func BuildString(s string) *Value {
	return &Value{stringReference, s}
}

func (v *Value) IsSymbol() bool { return v.ValueType == symbolReference }
func (v *Value) IsInteger() bool { return v.ValueType == integerReference }
func (v *Value) IsCons() bool { return v.ValueType == consReference }
func (v *Value) IsEmptyList() bool { return v.ValueType == emptyListReference }
func (v *Value) IsNativeFn() bool { return v.ValueType == nativeFnReference }
func (v *Value) IsString() bool { return v.ValueType == stringReference }

func (v *Value) IsList() bool {
	iter := v
	for iter.IsCons() {
		iter = iter.Cdr()
	}

	return iter.IsEmptyList()
}

func (lst *Value) ListLength() int {
	if !lst.IsList() {
		panic("Not a list")
	}

	res := 0
	iter := lst
	for iter.IsCons() {
		res += 1
		iter = iter.Cdr()
	}
	return res
}

func (sym *Value) SymbolName() string {
	if !sym.IsSymbol() {
		panic("Not a symbol")
	}

	return sym.Value.(string)
}

func (v *Value) IsLambdaSymbol() bool {
	return v.IsSymbol() && v.SymbolName() == "lambda"
}

func (n *Value) ToInt() int {
	if !n.IsInteger() {
		panic("Not an integer")
	}

	return n.Value.(int)
}

func (s *Value) ToStr() string {
	if !s.IsString() {
		panic("Not a string")
	}

	return s.Value.(string)
}

func (c *Value) Car() *Value {
	if c.IsEmptyList() {
		return c
	}

	if !c.IsCons() {
		panic("Not a cons")
	}

	return (*c.Value.(*Cons)).Car
}

func (c *Value) Cdr() *Value {
	if c.IsEmptyList() {
		return c
	}

	if !c.IsCons() {
		panic("Not a cons")
	}

	return (*c.Value.(*Cons)).Cdr
}

func (f *Value) NativeFn() (func(*Bindings, *Value) (*Value, error)) {
	return f.Value.(func(*Bindings, *Value) (*Value, error))
}


func Equal(a *Value, b *Value) bool {
	if a.ValueType != b.ValueType {
		return false
	}


	if a.IsEmptyList() {
		return true
	}

	if a.IsSymbol() {
		return a.SymbolName() == b.SymbolName()
	}

	if a.IsInteger() {
		return a.ToInt() == b.ToInt()
	}

	if a.IsCons() {
		return Equal(a.Car(), b.Car()) && Equal(a.Cdr(), b.Cdr())
	}

	if a.IsNativeFn() {
		return a.Value == b.Value
	}

	if a.IsString() {
		return a.Value == b.Value
	}

	panic("unexpected value type")
}
