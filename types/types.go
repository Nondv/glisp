package types

type Symbol string

const (
	symbolReference    = "sym"
	integerReference   = "int"
	consReference      = "cons"
	emptyListReference = "()"
)

type Value struct {
	ValueType string
	Value     any
}

type Cons struct {
	Car Value
	Cdr Value
}

func BuildSymbol(name string) Value {
	return Value{symbolReference, &name}
}

func BuildInteger(n int) Value {
	return Value{integerReference, n}
}

func BuildCons(car Value, cdr Value) Value {
	return Value{consReference, &Cons{car, cdr}}
}

func BuildEmptyList() Value {
	return Value{emptyListReference, nil}
}

func IsSymbol(v Value) bool { return v.ValueType == symbolReference }
func IsInteger(v Value) bool { return v.ValueType == integerReference }
func IsCons(v Value) bool { return v.ValueType == consReference }
func IsEmptyList(v Value) bool { return v.ValueType == emptyListReference }
