package expr

import "exprtree/value"

type Variable struct {
	Expr
	name string
}

func NewVariable(name string) *Variable {
	return &Variable{
		name: name,
	}
}

func (v *Variable) Name() string {
	return v.name
}

func (v *Variable) Eval() (value.Value, bool) {
	return nil, false
}

func (v *Variable) Equals(other any) bool {
	otherVar, ok := other.(*Variable)
	if !ok {
		return false
	}
	return v.name == otherVar.name
}
