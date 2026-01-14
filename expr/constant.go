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

func (c *Constant) Eval() (value.Value, bool) {
	return c.value, true
}

func (c *Constant) Equals(other any) bool {
	otherConst, ok := other.(*Constant)
	if !ok {
		return false
	}
	return c.value.Equals(otherConst.value)
}
