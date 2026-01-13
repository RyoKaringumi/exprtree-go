package prop

import (
	"exprtree/expr"
	"exprtree/value"
)

type Equal struct {
	Proposition
	left, right expr.Expr
}

func NewEqual(left, right expr.Expr) *Equal {
	return &Equal{
		left:  left,
		right: right,
	}
}

func (e *Equal) Left() expr.Expr {
	return e.left
}

func (e *Equal) Right() expr.Expr {
	return e.right
}

func (e *Equal) Eval() (value.Value, bool) {
	leftVal, ok := e.left.Eval()
	if !ok {
		return nil, false
	}
	rightVal, ok := e.right.Eval()
	if !ok {
		return nil, false
	}

	isEqual := leftVal == rightVal
	return value.NewBoolValue(isEqual), true
}
