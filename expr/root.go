package expr

import (
	"exprtree/value"
	"math"
)

// n-th root
type NthRoot struct {
	Expr
	radicand Expr
	degree   Expr
}

func NewNthRoot(radicand, degree Expr) *NthRoot {
	if radicand == nil {
		panic("radicand is nil")
	}
	if degree == nil {
		panic("degree is nil")
	}
	return &NthRoot{
		radicand: radicand,
		degree:   degree,
	}
}

func NewSqrt(radicand Expr) *NthRoot {
	return NewNthRoot(radicand, NewConstant(value.NewRealValue(2.0)))
}

func (n *NthRoot) Radicand() Expr {
	return n.radicand
}

func (n *NthRoot) Degree() Expr {
	return n.degree
}

func (n *NthRoot) Eval() (value.Value, bool) {
	leftVal, ok1 := n.radicand.Eval()
	rightVal, ok2 := n.degree.Eval()
	if !ok1 || !ok2 {
		return nil, false
	}
	leftReal, ok1 := leftVal.(*value.RealValue)
	rightReal, ok2 := rightVal.(*value.RealValue)
	if !ok1 || !ok2 || rightReal.Float64() == 0 {
		return nil, false
	}
	result := math.Pow(leftReal.Float64(), 1.0/rightReal.Float64())
	return value.NewRealValue(result), true
}
