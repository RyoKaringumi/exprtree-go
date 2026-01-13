package expr

import "exprtree/value"

type Mul struct {
	Binary
	left  Expr
	right Expr
}

func NewMul(left, right Expr) *Mul {
	if left == nil || right == nil {
		panic("left and right expressions must not be nil")
	}
	return &Mul{
		left:  left,
		right: right,
	}
}

func (m *Mul) Left() Expr {
	return m.left
}

func (m *Mul) Right() Expr {
	return m.right
}

func (m *Mul) Eval() (value.Value, bool) {
	leftVal, ok := m.left.Eval()
	if !ok {
		return nil, false
	}
	rightVal, ok := m.right.Eval()
	if !ok {
		return nil, false
	}

	leftReal, ok1 := leftVal.(*value.RealValue)
	rightReal, ok2 := rightVal.(*value.RealValue)
	if !ok1 || !ok2 {
		return nil, false
	}

	result := leftReal.Float64() * rightReal.Float64()
	return value.NewRealValue(result), true
}
