package expr

import "exprtree/value"

type Constant struct {
	Expr
	value value.Value
}

func NewConstant(value value.Value) *Constant {
	return &Constant{
		value: value,
	}
}

func (c *Constant) Value() value.Value {
	return c.value
}
