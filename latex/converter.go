package latex

import (
	"exprtree/expr"
	"exprtree/prop"
	"exprtree/value"
	"fmt"
)

// Converter converts LaTeX AST to Expression tree or Proposition
type Converter struct {
	errors []string
}

// NewConverter creates a new Converter instance
func NewConverter() *Converter {
	return &Converter{
		errors: []string{},
	}
}

// Convert converts a LatexNode to an Expression or Proposition
// Returns interface{} to support both expr.Expr and expr.Proposition types
func (c *Converter) Convert(node LatexNode) (interface{}, error) {
	if node == nil {
		return nil, fmt.Errorf("cannot convert nil node")
	}

	switch n := node.(type) {
	case *NumberNode:
		return c.convertNumber(n), nil
	case *VariableNode:
		return c.convertVariable(n), nil
	case *BinaryOpNode:
		return c.convertBinaryOp(n)
	case *EqualNode:
		return c.convertEqual(n)
	case *GroupNode:
		return c.convertGroup(n)
	case *CommandNode:
		return c.convertCommand(n)
	default:
		return nil, fmt.Errorf("unknown node type: %T", node)
	}
}

// convertNumber converts a NumberNode to a Constant
func (c *Converter) convertNumber(node *NumberNode) expr.Expr {
	return expr.NewConstant(value.NewRealValue(node.Value))
}

// convertVariable converts a VariableNode to a Variable
func (c *Converter) convertVariable(node *VariableNode) expr.Expr {
	return expr.NewVariable(node.Name)
}

// convertBinaryOp converts a BinaryOpNode to the appropriate Expression
func (c *Converter) convertBinaryOp(node *BinaryOpNode) (interface{}, error) {
	// Convert left and right children
	leftResult, err := c.Convert(node.Left)
	if err != nil {
		return nil, fmt.Errorf("failed to convert left operand: %w", err)
	}

	rightResult, err := c.Convert(node.Right)
	if err != nil {
		return nil, fmt.Errorf("failed to convert right operand: %w", err)
	}

	// Binary operations require Expression operands
	left, ok := leftResult.(expr.Expr)
	if !ok {
		return nil, fmt.Errorf("left operand must be an Expression, got %T", leftResult)
	}

	right, ok := rightResult.(expr.Expr)
	if !ok {
		return nil, fmt.Errorf("right operand must be an Expression, got %T", rightResult)
	}

	// Create appropriate Expression based on operator
	switch node.Operator.Type {
	case PLUS:
		return expr.NewAdd(left, right), nil
	case MINUS:
		return expr.NewSub(left, right), nil
	case MULTIPLY:
		return expr.NewMul(left, right), nil
	case DIVIDE:
		return expr.NewDiv(left, right), nil
	case CARET:
		return expr.NewPower(left, right), nil
	default:
		return nil, fmt.Errorf("unknown operator: %v", node.Operator.Type)
	}
}

// convertEqual converts an EqualNode to Equal or And proposition
// Handles chained equality: a = b = c becomes And(Eq(a,b), Eq(b,c))
func (c *Converter) convertEqual(node *EqualNode) (interface{}, error) {
	// Check if left side is also an EqualNode (chained equality)
	if leftEqualNode, ok := node.Left.(*EqualNode); ok {
		// This is a chained equality: (a = b) = c
		// Convert to: And(Eq(a, b), Eq(b, c))

		// Recursively convert left side (may produce And or Equal)
		leftResult, err := c.convertEqual(leftEqualNode)
		if err != nil {
			return nil, fmt.Errorf("failed to convert left equal: %w", err)
		}

		// Get the middle expression from the left equal node's right side
		middle, err := c.Convert(leftEqualNode.Right)
		if err != nil {
			return nil, fmt.Errorf("failed to convert middle operand: %w", err)
		}

		// Convert the rightmost expression
		right, err := c.Convert(node.Right)
		if err != nil {
			return nil, fmt.Errorf("failed to convert right operand: %w", err)
		}

		// Ensure middle and right are Expressions for Equal
		middleExpr, ok := middle.(expr.Expr)
		if !ok {
			return nil, fmt.Errorf("middle operand must be an Expression, got %T", middle)
		}

		rightExpr, ok := right.(expr.Expr)
		if !ok {
			return nil, fmt.Errorf("right operand must be an Expression, got %T", right)
		}

		// Create new Equal(middle, right)
		newEqual := prop.NewEqual(middleExpr, rightExpr)

		// Convert leftResult to Proposition
		leftProp, ok := leftResult.(prop.Proposition)
		if !ok {
			return nil, fmt.Errorf("left result must be a Proposition, got %T", leftResult)
		}

		// Return And(leftResult, Equal(middle, right))
		return prop.NewAnd(leftProp, newEqual), nil
	}

	// Not a chained equality, convert as simple Equal
	left, err := c.Convert(node.Left)
	if err != nil {
		return nil, fmt.Errorf("failed to convert left operand: %w", err)
	}

	right, err := c.Convert(node.Right)
	if err != nil {
		return nil, fmt.Errorf("failed to convert right operand: %w", err)
	}

	// Ensure left and right are Expressions
	leftExpr, ok := left.(expr.Expr)
	if !ok {
		return nil, fmt.Errorf("left operand must be an Expression, got %T", left)
	}

	rightExpr, ok := right.(expr.Expr)
	if !ok {
		return nil, fmt.Errorf("right operand must be an Expression, got %T", right)
	}

	return prop.NewEqual(leftExpr, rightExpr), nil
}

// convertGroup converts a GroupNode by converting its inner expression
func (c *Converter) convertGroup(node *GroupNode) (interface{}, error) {
	// Groups are just for parsing precedence, we don't need them in the Expression tree
	return c.Convert(node.Inner)
}

// convertCommand converts a CommandNode to the appropriate Expression
func (c *Converter) convertCommand(node *CommandNode) (interface{}, error) {
	switch node.Name {
	case "sqrt":
		return c.convertSqrt(node)
	default:
		return nil, fmt.Errorf("unknown command: \\%s", node.Name)
	}
}

// convertSqrt converts \sqrt command to SqrtExpression
func (c *Converter) convertSqrt(node *CommandNode) (interface{}, error) {
	// Convert the argument (radicand)
	argumentResult, err := c.Convert(node.Argument)
	if err != nil {
		return nil, fmt.Errorf("failed to convert sqrt argument: %w", err)
	}

	// Sqrt requires Expression argument
	argument, ok := argumentResult.(expr.Expr)
	if !ok {
		return nil, fmt.Errorf("sqrt argument must be an Expression, got %T", argumentResult)
	}

	// Handle optional root degree [n]
	if node.Optional != nil {
		optionalResult, err := c.Convert(node.Optional)
		if err != nil {
			return nil, fmt.Errorf("failed to convert sqrt optional argument: %w", err)
		}

		// Extract numeric value from optional argument
		optionalExpr, ok := optionalResult.(expr.Expr)
		if !ok {
			return nil, fmt.Errorf("sqrt optional argument must be an Expression, got %T", optionalResult)
		}

		constant, ok := optionalExpr.(*expr.Constant)
		if !ok {
			return nil, fmt.Errorf("sqrt optional argument must be a constant number")
		}

		return expr.NewNthRoot(argument, constant), nil
	}

	// Default to square root (N=2)
	return expr.NewSqrt(argument), nil
}

// Errors returns the list of conversion errors
func (c *Converter) Errors() []string {
	return c.errors
}
