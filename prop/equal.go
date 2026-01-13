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

	switch leftVal.Kind() {
	case value.RealKind:
		leftReal, ok1 := leftVal.(*value.RealValue)
		rightReal, ok2 := rightVal.(*value.RealValue)
		if !ok1 || !ok2 {
			return nil, false
		}
		result := leftReal.Float64() == rightReal.Float64()
		return value.NewBoolValue(result), true
	default:
		return nil, false
	}
}
