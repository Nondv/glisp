package types

type Bindings struct {
	SymbolName string
	Value      *Value
	Next       *Bindings
}

func (b *Bindings) Lookup(sym *Value) (*Value, bool) {
	name := sym.SymbolName()
	next := b
	for next != nil {
		if name == next.SymbolName {
			return next.Value, true
		}
		next = next.Next
	}

	return nil, false
}

func (b *Bindings) Assoc(sym *Value, val *Value) *Bindings {
	return &Bindings{sym.SymbolName(), val, b}
}

func (b *Bindings) AssocSym(sym string, val *Value) *Bindings {
	return b.Assoc(BuildSymbol(sym), val)
}
