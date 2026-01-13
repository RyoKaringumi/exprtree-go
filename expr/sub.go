package expr

import "exprtree/value"

type Sub struct {
	Binary
	left, right Expr
}

func NewSub(left, right Expr) *Sub {
	if left == nil || right == nil {
		panic("left or right is nil")
	}
	return &Sub{
		left:  left,
		right: right,
	}
}

func (s *Sub) Left() Expr {
	return s.left
}

func (s *Sub) Right() Expr {
	return s.right
}

func (s *Sub) Eval() (value.Value, bool) {
	leftVal, ok := s.left.Eval()
	if !ok {
		return nil, false
	}
	rightVal, ok := s.right.Eval()
	if !ok {
		return nil, false
	}

	leftReal, ok1 := leftVal.(*value.RealValue)
	rightReal, ok2 := rightVal.(*value.RealValue)
	if !ok1 || !ok2 {
		return nil, false
	}

	result := leftReal.Float64() - rightReal.Float64()
	return value.NewRealValue(result), true
}
