package latex

import (
	"exprtree/expr"
	"fmt"
)

// Converter converts LaTeX AST to Expression tree
type Converter struct {
	errors []string
}

// NewConverter creates a new Converter instance
func NewConverter() *Converter {
	return &Converter{
		errors: []string{},
	}
}

// Convert converts a LatexNode to an Expression
func (c *Converter) Convert(node LatexNode) (expr.Expression, error) {
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
func (c *Converter) convertNumber(node *NumberNode) expr.Expression {
	return &expr.Constant{
		Value: expr.NumberValue{Value: node.Value},
	}
}

// convertVariable converts a VariableNode to a Variable
func (c *Converter) convertVariable(node *VariableNode) expr.Expression {
	return &expr.Variable{
		Name: node.Name,
	}
}

// convertBinaryOp converts a BinaryOpNode to the appropriate Expression
func (c *Converter) convertBinaryOp(node *BinaryOpNode) (expr.Expression, error) {
	// Convert left and right children
	left, err := c.Convert(node.Left)
	if err != nil {
		return nil, fmt.Errorf("failed to convert left operand: %w", err)
	}

	right, err := c.Convert(node.Right)
	if err != nil {
		return nil, fmt.Errorf("failed to convert right operand: %w", err)
	}

	// Create appropriate Expression based on operator
	switch node.Operator.Type {
	case PLUS:
		return expr.NewAddExpression(left, right), nil
	case MINUS:
		return expr.NewSubtractExpression(left, right), nil
	case MULTIPLY:
		return expr.NewMultiplyExpression(left, right), nil
	case DIVIDE:
		return expr.NewDivideExpression(left, right), nil
	case CARET:
		return expr.NewPowerExpression(left, right), nil
	default:
		return nil, fmt.Errorf("unknown operator: %v", node.Operator.Type)
	}
}

// convertEqual converts an EqualNode to EqualExpression
func (c *Converter) convertEqual(node *EqualNode) (expr.Expression, error) {
	// Convert left and right children
	left, err := c.Convert(node.Left)
	if err != nil {
		return nil, fmt.Errorf("failed to convert left operand: %w", err)
	}

	right, err := c.Convert(node.Right)
	if err != nil {
		return nil, fmt.Errorf("failed to convert right operand: %w", err)
	}

	return expr.NewEqualExpression(left, right), nil
}

// convertGroup converts a GroupNode by converting its inner expression
func (c *Converter) convertGroup(node *GroupNode) (expr.Expression, error) {
	// Groups are just for parsing precedence, we don't need them in the Expression tree
	return c.Convert(node.Inner)
}

// convertCommand converts a CommandNode to the appropriate Expression
func (c *Converter) convertCommand(node *CommandNode) (expr.Expression, error) {
	switch node.Name {
	case "sqrt":
		return c.convertSqrt(node)
	default:
		return nil, fmt.Errorf("unknown command: \\%s", node.Name)
	}
}

// convertSqrt converts \sqrt command to SqrtExpression
func (c *Converter) convertSqrt(node *CommandNode) (expr.Expression, error) {
	// Convert the argument (radicand)
	argument, err := c.Convert(node.Argument)
	if err != nil {
		return nil, fmt.Errorf("failed to convert sqrt argument: %w", err)
	}

	// Handle optional root degree [n]
	if node.Optional != nil {
		optionalExpr, err := c.Convert(node.Optional)
		if err != nil {
			return nil, fmt.Errorf("failed to convert sqrt optional argument: %w", err)
		}

		// Extract numeric value from optional argument
		constant, ok := optionalExpr.(*expr.Constant)
		if !ok {
			return nil, fmt.Errorf("sqrt optional argument must be a constant number")
		}

		return expr.NewNthRootExpression(argument, constant.Value.Value), nil
	}

	// Default to square root (N=2)
	return expr.NewSqrtExpression(argument), nil
}

// Errors returns the list of conversion errors
func (c *Converter) Errors() []string {
	return c.errors
}
