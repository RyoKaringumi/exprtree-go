package expr

import "exprtree/value"

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
