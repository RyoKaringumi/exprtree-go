package expr

import "exprtree/value"

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
	return nil, false
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
