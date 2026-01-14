package expr

import "exprtree/value"

type Add struct {
	Binary
	left  Expr
	right Expr
}

func NewAdd(left Expr, right Expr) *Add {
	if left == nil || right == nil {
		panic("left and right expressions must not be nil")
	}
	return &Add{
		left:  left,
		right: right,
	}
}

func (a *Add) Left() Expr {
	return a.left
}

func (a *Add) Right() Expr {
	return a.right
}

func (a *Add) Eval() (value.Value, bool) {
	leftVal, ok := a.left.Eval()
	if !ok {
		return nil, false
	}
	rightVal, ok := a.right.Eval()
	if !ok {
		return nil, false
	}

	leftReal, ok1 := leftVal.(*value.RealValue)
	rightReal, ok2 := rightVal.(*value.RealValue)
	if !ok1 || !ok2 {
		return nil, false
	}

	result := leftReal.Float64() + rightReal.Float64()
	return value.NewRealValue(result), true
}

func (a *Add) Equals(other any) bool {
	otherAdd, ok := other.(*Add)
	if !ok {
		return false
	}
	return a.left.Equals(otherAdd.left) && a.right.Equals(otherAdd.right)
}
