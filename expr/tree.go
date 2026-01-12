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

type BoolValue struct {
	Value bool
}

func (b *BoolValue) Eval() (ExpressValue, bool) {
	return b, true
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

type Add struct {
	BinaryExpression
}

func NewAdd(left, right Expression) *Add {
	return &Add{
		BinaryExpression: BinaryExpression{
			Left:  left,
			Right: right,
		},
	}
}

type Subtract struct {
	BinaryExpression
}

func NewSubtract(left, right Expression) *Subtract {
	return &Subtract{
		BinaryExpression: BinaryExpression{
			Left:  left,
			Right: right,
		},
	}
}

type Multiply struct {
	BinaryExpression
}

func NewMultiply(left, right Expression) *Multiply {
	return &Multiply{
		BinaryExpression: BinaryExpression{
			Left:  left,
			Right: right,
		},
	}
}

type Divide struct {
	BinaryExpression
}

func NewDivide(left, right Expression) *Divide {
	return &Divide{
		BinaryExpression: BinaryExpression{
			Left:  left,
			Right: right,
		},
	}
}

type Proposition interface {
}

type Equal struct {
	Proposition
	BinaryExpression
}

func NewEqual(left, right Expression) *Equal {
	return &Equal{
		BinaryExpression: BinaryExpression{
			Left:  left,
			Right: right,
		},
	}
}

type And struct {
	Proposition
	Left  Proposition
	Right Proposition
}

func NewAnd(left, right Proposition) *And {
	return &And{
		Left:  left,
		Right: right,
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

func (e *Add) Eval() (ExpressValue, bool) {
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

func (e *Subtract) Eval() (ExpressValue, bool) {
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

func (e *Multiply) Eval() (ExpressValue, bool) {
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

func (e *Divide) Eval() (ExpressValue, bool) {
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

func (e *Equal) Eval() (ExpressValue, bool) {
	leftVal, leftOk := e.Left.Eval()
	rightVal, rightOk := e.Right.Eval()

	if !leftOk || !rightOk {
		return nil, false
	}

	// NumberValue同士の比較
	if leftNum, ok := leftVal.(*NumberValue); ok {
		if rightNum, ok := rightVal.(*NumberValue); ok {
			// 浮動小数点誤差を考慮した比較
			const epsilon = 1e-9
			equal := math.Abs(leftNum.Value-rightNum.Value) < epsilon
			return &BoolValue{Value: equal}, true
		}
	}

	// BoolValue同士の比較
	if leftBool, ok := leftVal.(*BoolValue); ok {
		if rightBool, ok := rightVal.(*BoolValue); ok {
			return &BoolValue{Value: leftBool.Value == rightBool.Value}, true
		}
	}

	// 型が異なる場合は評価失敗
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
	case *Add:
		if e, ok := expr.(*Add); ok {
			return patternMatchHelper(p.Left, e.Left, bindings) &&
				patternMatchHelper(p.Right, e.Right, bindings)
		}
	case *Subtract:
		if e, ok := expr.(*Subtract); ok {
			return patternMatchHelper(p.Left, e.Left, bindings) &&
				patternMatchHelper(p.Right, e.Right, bindings)
		}
	case *Multiply:
		if e, ok := expr.(*Multiply); ok {
			return patternMatchHelper(p.Left, e.Left, bindings) &&
				patternMatchHelper(p.Right, e.Right, bindings)
		}
	case *Divide:
		if e, ok := expr.(*Divide); ok {
			return patternMatchHelper(p.Left, e.Left, bindings) &&
				patternMatchHelper(p.Right, e.Right, bindings)
		}
	case *Power:
		if e, ok := expr.(*Power); ok {
			return patternMatchHelper(p.Left, e.Left, bindings) &&
				patternMatchHelper(p.Right, e.Right, bindings)
		}
	case *Equal:
		if e, ok := expr.(*Equal); ok {
			return patternMatchHelper(p.Left, e.Left, bindings) &&
				patternMatchHelper(p.Right, e.Right, bindings)
		}
	case *Sqrt:
		if e, ok := expr.(*Sqrt); ok {
			// Root degree must match exactly
			if p.N != e.N {
				return false
			}
			return patternMatchHelper(p.Operand, e.Operand, bindings)
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
	case *Add:
		if e2, ok := expr2.(*Add); ok {
			return expressionsEqual(e1.Left, e2.Left) && expressionsEqual(e1.Right, e2.Right)
		}
	case *Subtract:
		if e2, ok := expr2.(*Subtract); ok {
			return expressionsEqual(e1.Left, e2.Left) && expressionsEqual(e1.Right, e2.Right)
		}
	case *Multiply:
		if e2, ok := expr2.(*Multiply); ok {
			return expressionsEqual(e1.Left, e2.Left) && expressionsEqual(e1.Right, e2.Right)
		}
	case *Divide:
		if e2, ok := expr2.(*Divide); ok {
			return expressionsEqual(e1.Left, e2.Left) && expressionsEqual(e1.Right, e2.Right)
		}
	case *Power:
		if e2, ok := expr2.(*Power); ok {
			return expressionsEqual(e1.Left, e2.Left) && expressionsEqual(e1.Right, e2.Right)
		}
	case *Equal:
		if e2, ok := expr2.(*Equal); ok {
			return expressionsEqual(e1.Left, e2.Left) && expressionsEqual(e1.Right, e2.Right)
		}
	case *Sqrt:
		if e2, ok := expr2.(*Sqrt); ok {
			return e1.N == e2.N && expressionsEqual(e1.Operand, e2.Operand)
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
	case *Add:
		return NewAdd(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *Subtract:
		return NewSubtract(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *Multiply:
		return NewMultiply(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *Divide:
		return NewDivide(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *Power:
		return NewPower(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *Equal:
		return NewEqual(
			Substitute(e.Left, bindings),
			Substitute(e.Right, bindings),
		)
	case *Sqrt:
		return NewNthRoot(
			Substitute(e.Operand, bindings),
			e.N,
		)
	}

	// Unknown expression type, return as-is
	return expr
}

type Power struct {
	BinaryExpression
}

func NewPower(base, exponent Expression) *Power {
	return &Power{
		BinaryExpression: BinaryExpression{
			Left:  base,
			Right: exponent,
		},
	}
}

func (e *Power) Eval() (ExpressValue, bool) {
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

type Sqrt struct {
	Expression
	Operand Expression
	N       float64 // Root degree (e.g., 2 for square root, 3 for cube root)
}

func NewSqrt(operand Expression) *Sqrt {
	return NewNthRoot(operand, 2)
}

func NewNthRoot(operand Expression, n float64) *Sqrt {
	return &Sqrt{
		Operand: operand,
		N:       n,
	}
}

func (e *Sqrt) Eval() (ExpressValue, bool) {
	operandVal, operandOk := e.Operand.Eval()

	if !operandOk {
		return nil, false
	}
	if operandNum, ok := operandVal.(*NumberValue); ok {
		// For even roots, operand must be non-negative
		// For odd roots, negative values are allowed
		isEvenRoot := int(e.N)%2 == 0
		if isEvenRoot && operandNum.Value < 0 {
			return nil, false
		}

		// Calculate nth root: x^(1/n)
		// For negative values with odd roots, handle sign separately
		if operandNum.Value < 0 {
			result := -math.Pow(-operandNum.Value, 1.0/e.N)
			return &NumberValue{Value: result}, true
		}

		return &NumberValue{Value: math.Pow(operandNum.Value, 1.0/e.N)}, true
	}
	return nil, false
}

func (e *Sqrt) Children() []Expression {
	return []Expression{e.Operand}
}
