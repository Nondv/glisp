package types

type Bindings struct {
	SymbolName string
	Value      Value
	Next       *Bindings
}

func (b *Bindings) Lookup(sym Value) (Value, bool) {
	if !sym.IsSymbol() {
		panic("not a symbol")
	}

	name := *sym.Value.(*string)

	next := b
	for next != nil {
		if name == next.SymbolName {
			return next.Value, true
		}
		next = next.Next
	}

	return Value{}, false
}

func (b *Bindings) Assoc(sym Value, val Value) *Bindings {
	if !sym.IsSymbol() {
		panic("not a symbol")
	}

	name := *sym.Value.(*string)
	return &Bindings{name, val, b}
}
