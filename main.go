package main

type ExpressValue interface {
}

type NumberValue struct {
	ExpressValue
	Value float64
}

type Expression interface {
	eval() ExpressValue
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

func (e *AddExpression) eval() ExpressValue {
	leftVal := e.Left.eval()
	rightVal := e.Right.eval()

	if leftNum, ok := leftVal.(*NumberValue); ok {
		if rightNum, ok := rightVal.(*NumberValue); ok {
			return &NumberValue{Value: leftNum.Value + rightNum.Value}
		}
	}
	return nil
}

func (e *SubtractExpression) eval() ExpressValue {
	leftVal := e.Left.eval()
	rightVal := e.Right.eval()

	if leftNum, ok := leftVal.(*NumberValue); ok {
		if rightNum, ok := rightVal.(*NumberValue); ok {
			return &NumberValue{Value: leftNum.Value - rightNum.Value}
		}
	}
	return nil
}

func (e *MultiplyExpression) eval() ExpressValue {
	leftVal := e.Left.eval()
	rightVal := e.Right.eval()

	if leftNum, ok := leftVal.(*NumberValue); ok {
		if rightNum, ok := rightVal.(*NumberValue); ok {
			return &NumberValue{Value: leftNum.Value * rightNum.Value}
		}
	}
	return nil
}

func (e *DivideExpression) eval() ExpressValue {
	leftVal := e.Left.eval()
	rightVal := e.Right.eval()

	if leftNum, ok := leftVal.(*NumberValue); ok {
		if rightNum, ok := rightVal.(*NumberValue); ok {
			if rightNum.Value != 0 {
				return &NumberValue{Value: leftNum.Value / rightNum.Value}
			}
		}
	}
	return nil
}

func main() {
	// This is a placeholder for the main function.
}
