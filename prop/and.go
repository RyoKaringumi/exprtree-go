package prop

import "exprtree/value"

type And struct {
	Proposition
	left, right Proposition
}

func NewAnd(left, right Proposition) *And {
	return &And{
		left:  left,
		right: right,
	}
}

func (a *And) Left() Proposition {
	return a.left
}

func (a *And) Right() Proposition {
	return a.right
}

func (a *And) Eval() (value.Value, bool) {
	leftVal, ok := a.left.Eval()
	if !ok {
		return nil, false
	}
	rightVal, ok := a.right.Eval()
	if !ok {
		return nil, false
	}

	if leftVal.Kind() != value.BoolKind || rightVal.Kind() != value.BoolKind {
		return nil, false
	}

	result := leftVal.(*value.BoolValue).Bool() && rightVal.(*value.BoolValue).Bool()
	return value.NewBoolValue(result), true
}
