package expr

import (
	"exprtree/ast"
	"exprtree/value"
	"math"
)

type Power struct {
	Binary
	base     Expr
	exponent Expr
}

func NewPower(base, exponent Expr) *Power {
	if base == nil {
		panic("base is nil")
	}
	if exponent == nil {
		panic("exponent is nil")
	}
	return &Power{
		base:     base,
		exponent: exponent,
	}
}

func (p *Power) Eval() (value.Value, bool) {
	baseVal, ok := p.base.Eval()
	if !ok {
		return nil, false
	}
	exponentVal, ok := p.exponent.Eval()
	if !ok {
		return nil, false
	}

	baseReal, ok1 := baseVal.(*value.RealValue)
	exponentReal, ok2 := exponentVal.(*value.RealValue)
	if !ok1 || !ok2 {
		return nil, false
	}

	result := math.Pow(baseReal.Float64(), exponentReal.Float64())
	return value.NewRealValue(result), true
}

func (p *Power) Base() Expr {
	return p.base
}

func (p *Power) Exponent() Expr {
	return p.exponent
}

func (p *Power) Left() Expr {
	return p.base
}

func (p *Power) Right() Expr {
	return p.exponent
}

func (p *Power) Equals(other any) bool {
	otherPower, ok := other.(*Power)
	if !ok {
		return false
	}
	return p.base.Equals(otherPower.base) && p.exponent.Equals(otherPower.exponent)
}

func (p *Power) Children() []ast.HasChildren {
	return []ast.HasChildren{p.base, p.exponent}
}
