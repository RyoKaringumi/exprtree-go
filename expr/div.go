package expr

import (
	"exprtree/ast"
	"exprtree/value"
)

type Div struct {
	Binary
	left  Expr
	right Expr
}

func NewDiv(left, right Expr) *Div {
	if left == nil || right == nil {
		panic("left and right expressions must not be nil")
	}
	return &Div{
		left:  left,
		right: right,
	}
}

func (d *Div) Left() Expr {
	return d.left
}

func (d *Div) Right() Expr {
	return d.right
}

func (d *Div) Eval() (value.Value, bool) {
	leftVal, ok1 := d.left.Eval()
	rightVal, ok2 := d.right.Eval()
	if !ok1 || !ok2 {
		return nil, false
	}
	leftReal, ok1 := leftVal.(*value.RealValue)
	rightReal, ok2 := rightVal.(*value.RealValue)
	if !ok1 || !ok2 || rightReal.Float64() == 0 {
		return nil, false
	}
	result := leftReal.Float64() / rightReal.Float64()
	return value.NewRealValue(result), true
}

func (d *Div) Equals(other any) bool {
	otherDiv, ok := other.(*Div)
	if !ok {
		return false
	}
	return d.left.Equals(otherDiv.left) && d.right.Equals(otherDiv.right)
}

func (d *Div) Children() []ast.HasChildren {
	return []ast.HasChildren{d.left, d.right}
}
