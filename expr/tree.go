package expr

type ExpressValue interface {
}

type NumberValue struct {
	Value float64
}

func (n *NumberValue) Eval() (ExpressValue, bool) {
	return n, true
}

type Expression interface {
	Children() []Expression
	Eval() (ExpressValue, bool)
}

type BinaryExpression struct {
	Expression
	Left  Expression
	Right Expression
}

func (b *BinaryExpression) Children() []Expression {
	return []Expression{b.Left, b.Right}
}

type AddExpression struct {
	BinaryExpression
}

func NewAddExpression(left, right Expression) *AddExpression {
	return &AddExpression{
		BinaryExpression: BinaryExpression{
			Left:  left,
			Right: right,
		},
	}
}

type SubtractExpression struct {
	BinaryExpression
}

func NewSubtractExpression(left, right Expression) *SubtractExpression {
	return &SubtractExpression{
		BinaryExpression: BinaryExpression{
			Left:  left,
			Right: right,
		},
	}
}

type MultiplyExpression struct {
	BinaryExpression
}

func NewMultiplyExpression(left, right Expression) *MultiplyExpression {
	return &MultiplyExpression{
		BinaryExpression: BinaryExpression{
			Left:  left,
			Right: right,
		},
	}
}

type DivideExpression struct {
	BinaryExpression
}

func NewDivideExpression(left, right Expression) *DivideExpression {
	return &DivideExpression{
		BinaryExpression: BinaryExpression{
			Left:  left,
			Right: right,
		},
	}
}

type Constant struct {
	Expression
	Value NumberValue
}

type Variable struct {
	Expression
	Name string
}

func (c *Constant) Children() []Expression {
	return []Expression{}
}

func (v *Variable) Children() []Expression {
	return []Expression{}
}

func (e *AddExpression) Eval() (ExpressValue, bool) {
	leftVal, leftOk := e.Left.Eval()
	rightVal, rightOk := e.Right.Eval()

	if !leftOk || !rightOk {
		return nil, false
	}

	if leftNum, ok := leftVal.(*NumberValue); ok {
		if rightNum, ok := rightVal.(*NumberValue); ok {
			return &NumberValue{Value: leftNum.Value + rightNum.Value}, true
		}
	}
	return nil, false
}

func (e *SubtractExpression) Eval() (ExpressValue, bool) {
	leftVal, leftOk := e.Left.Eval()
	rightVal, rightOk := e.Right.Eval()

	if !leftOk || !rightOk {
		return nil, false
	}

	if leftNum, ok := leftVal.(*NumberValue); ok {
		if rightNum, ok := rightVal.(*NumberValue); ok {
			return &NumberValue{Value: leftNum.Value - rightNum.Value}, true
		}
	}
	return nil, false
}

func (e *MultiplyExpression) Eval() (ExpressValue, bool) {
	leftVal, leftOk := e.Left.Eval()
	rightVal, rightOk := e.Right.Eval()

	if !leftOk || !rightOk {
		return nil, false
	}

	if leftNum, ok := leftVal.(*NumberValue); ok {
		if rightNum, ok := rightVal.(*NumberValue); ok {
			return &NumberValue{Value: leftNum.Value * rightNum.Value}, true
		}
	}
	return nil, false
}

func (e *DivideExpression) Eval() (ExpressValue, bool) {
	leftVal, leftOk := e.Left.Eval()
	rightVal, rightOk := e.Right.Eval()

	if !leftOk || !rightOk {
		return nil, false
	}

	if leftNum, ok := leftVal.(*NumberValue); ok {
		if rightNum, ok := rightVal.(*NumberValue); ok {
			if rightNum.Value != 0 {
				return &NumberValue{Value: leftNum.Value / rightNum.Value}, true
			}
		}
	}
	return nil, false
}

func (c *Constant) Eval() (ExpressValue, bool) {
	return &c.Value, true
}
