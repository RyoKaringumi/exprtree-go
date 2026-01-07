package expr

import "math"

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

func NewConstant(value float64) *Constant {
	return &Constant{
		Value: NumberValue{Value: value},
	}
}

type Variable struct {
	Expression
	Name string
}

func NewVariable(name string) *Variable {
	return &Variable{
		Name: name,
	}
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

func (v *Variable) Eval() (ExpressValue, bool) {
	// Variables cannot be evaluated without a value assignment context
	return nil, false
}

// PatternMatch attempts to match an expression against a pattern template.
// Variables in the pattern act as wildcards that can match any sub-expression.
// Returns a map of variable names to their matched expressions, and a success flag.
//
// Example:
//
//	pattern: x(y+z)  → x matches (x+3), y matches 1, z matches z
//	expr: (x+3)(1+z)
func PatternMatch(pattern Expression, expr Expression) (map[string]Expression, bool) {
	bindings := make(map[string]Expression)
	ok := patternMatchHelper(pattern, expr, bindings)
	return bindings, ok
}

func patternMatchHelper(pattern Expression, expr Expression, bindings map[string]Expression) bool {
	// If pattern is a variable, treat it as a wildcard
	if patternVar, ok := pattern.(*Variable); ok {
		// If this variable is already bound, check if expr matches the existing binding
		if existingExpr, exists := bindings[patternVar.Name]; exists {
			return expressionsEqual(existingExpr, expr)
		}
		// Create new binding
		bindings[patternVar.Name] = expr
		return true
	}

	// If pattern is a constant, expr must be the same constant
	if patternConst, ok := pattern.(*Constant); ok {
		if exprConst, ok := expr.(*Constant); ok {
			return patternConst.Value.Value == exprConst.Value.Value
		}
		return false
	}

	// For binary expressions, both must be the same type and children must match
	switch p := pattern.(type) {
	case *AddExpression:
		if e, ok := expr.(*AddExpression); ok {
			return patternMatchHelper(p.Left, e.Left, bindings) &&
				patternMatchHelper(p.Right, e.Right, bindings)
		}
	case *SubtractExpression:
		if e, ok := expr.(*SubtractExpression); ok {
			return patternMatchHelper(p.Left, e.Left, bindings) &&
				patternMatchHelper(p.Right, e.Right, bindings)
		}
	case *MultiplyExpression:
		if e, ok := expr.(*MultiplyExpression); ok {
			return patternMatchHelper(p.Left, e.Left, bindings) &&
				patternMatchHelper(p.Right, e.Right, bindings)
		}
	case *DivideExpression:
		if e, ok := expr.(*DivideExpression); ok {
			return patternMatchHelper(p.Left, e.Left, bindings) &&
				patternMatchHelper(p.Right, e.Right, bindings)
		}
	}

	return false
}

// expressionsEqual checks if two expressions are structurally equal
func expressionsEqual(expr1, expr2 Expression) bool {
	// Check if both are constants
	if c1, ok := expr1.(*Constant); ok {
		if c2, ok := expr2.(*Constant); ok {
			return c1.Value.Value == c2.Value.Value
		}
		return false
	}

	// Check if both are variables with the same name
	if v1, ok := expr1.(*Variable); ok {
		if v2, ok := expr2.(*Variable); ok {
			return v1.Name == v2.Name
		}
		return false
	}

	// Check if both are the same type of binary expression
	switch e1 := expr1.(type) {
	case *AddExpression:
		if e2, ok := expr2.(*AddExpression); ok {
			return expressionsEqual(e1.Left, e2.Left) && expressionsEqual(e1.Right, e2.Right)
		}
	case *SubtractExpression:
		if e2, ok := expr2.(*SubtractExpression); ok {
			return expressionsEqual(e1.Left, e2.Left) && expressionsEqual(e1.Right, e2.Right)
		}
	case *MultiplyExpression:
		if e2, ok := expr2.(*MultiplyExpression); ok {
			return expressionsEqual(e1.Left, e2.Left) && expressionsEqual(e1.Right, e2.Right)
		}
	case *DivideExpression:
		if e2, ok := expr2.(*DivideExpression); ok {
			return expressionsEqual(e1.Left, e2.Left) && expressionsEqual(e1.Right, e2.Right)
		}
	}

	return false
}

// Substitute replaces variables in an expression with their bound expressions.
// This is the inverse operation of PatternMatch.
//
// Example:
//
//	expr: x(y+z)
//	bindings: {x: (x+3), y: 1, z: z}
//	result: (x+3)(1+z)
func Substitute(expr Expression, bindings map[string]Expression) Expression {
	// If expr is a variable, replace it if a binding exists
	if v, ok := expr.(*Variable); ok {
		if replacement, exists := bindings[v.Name]; exists {
			return replacement
		}
		// No binding, return the variable as-is
		return v
	}

	// If expr is a constant, return it as-is
	if c, ok := expr.(*Constant); ok {
		return c
	}

	// For binary expressions, recursively substitute in children
	switch e := expr.(type) {
	case *AddExpression:
		return NewAddExpression(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *SubtractExpression:
		return NewSubtractExpression(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *MultiplyExpression:
		return NewMultiplyExpression(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *DivideExpression:
		return NewDivideExpression(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *PowerExpression:
		return NewPowerExpression(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *SqrtExpression:
		return NewSqrtExpression(
			Substitute(e.Operand, bindings),
		)
	}

	// Unknown expression type, return as-is
	return expr
}

type PowerExpression struct {
	BinaryExpression
}

func NewPowerExpression(base, exponent Expression) *PowerExpression {
	return &PowerExpression{
		BinaryExpression: BinaryExpression{
			Left:  base,
			Right: exponent,
		},
	}
}

func (e *PowerExpression) Eval() (ExpressValue, bool) {
	baseVal, baseOk := e.Left.Eval()
	exponentVal, exponentOk := e.Right.Eval()

	if !baseOk || !exponentOk {
		return nil, false
	}

	if baseNum, ok := baseVal.(*NumberValue); ok {
		if exponentNum, ok := exponentVal.(*NumberValue); ok {
			return &NumberValue{Value: math.Pow(baseNum.Value, exponentNum.Value)}, true
		}
	}
	return nil, false
}

type SqrtExpression struct {
	Expression
	Operand Expression
}

func NewSqrtExpression(operand Expression) *SqrtExpression {
	return &SqrtExpression{
		Operand: operand,
	}
}

func (e *SqrtExpression) Eval() (ExpressValue, bool) {
	operandVal, operandOk := e.Operand.Eval()

	if !operandOk {
		return nil, false
	}
	if operandNum, ok := operandVal.(*NumberValue); ok {
		if operandNum.Value >= 0 {
			return &NumberValue{Value: math.Sqrt(operandNum.Value)}, true
		}
	}
	return nil, false
}

func (e *SqrtExpression) Children() []Expression {
	return []Expression{e.Operand}
}
