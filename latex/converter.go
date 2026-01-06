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
	case *GroupNode:
		return c.convertGroup(n)
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
	default:
		return nil, fmt.Errorf("unknown operator: %v", node.Operator.Type)
	}
}

// convertGroup converts a GroupNode by converting its inner expression
func (c *Converter) convertGroup(node *GroupNode) (expr.Expression, error) {
	// Groups are just for parsing precedence, we don't need them in the Expression tree
	return c.Convert(node.Inner)
}

// Errors returns the list of conversion errors
func (c *Converter) Errors() []string {
	return c.errors
}
