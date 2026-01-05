package main

type ExpressValue interface {
}

type NumberValue struct {
	Value float64
}

func (n *NumberValue) Eval() (ExpressValue, bool) {
	return n, true
}

type Expression interface {
	Eval() (ExpressValue, bool)
}

type AddExpression struct {
	Expression
	Left  Expression
	Right Expression
}

type SubtractExpression struct {
	Expression
	Left  Expression
	Right Expression
}

type MultiplyExpression struct {
	Expression
	Left  Expression
	Right Expression
}

type DivideExpression struct {
	Expression
	Left  Expression
	Right Expression
}

type Constant struct {
	Expression
	Value NumberValue
}

type Variable struct {
	Expression
	Name string
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

func (v *Variable) Eval() (ExpressValue, bool) {
	// Variable evaluation logic would go here.
	// For simplicity, returning nil, false.
	return nil, false
}
