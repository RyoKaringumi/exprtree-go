package latex

import (
	"exprtree/expr"
	"fmt"
	"strings"
)

// Renderer converts LaTeX AST to string representation
type Renderer struct {
}

// NewRenderer creates a new Renderer instance
func NewRenderer() *Renderer {
	return &Renderer{}
}

// Render converts a LatexNode to a string
func (r *Renderer) Render(node LatexNode) string {
	return r.renderNode(node, LOWEST, false)
}

// renderNode converts a LatexNode to a string with precedence handling
// parentPrec is the precedence of the parent operator
// isRightOperand indicates if this node is the right operand of a binary operator
func (r *Renderer) renderNode(node LatexNode, parentPrec int, isRightOperand bool) string {
	switch n := node.(type) {
	case *NumberNode:
		return r.renderNumber(n)
	case *BinaryOpNode:
		return r.renderBinaryOp(n, parentPrec, isRightOperand)
	case *GroupNode:
		return r.renderGroup(n)
	default:
		return ""
	}
}

// renderNumber converts a NumberNode to a string
func (r *Renderer) renderNumber(node *NumberNode) string {
	// Use %g to format the number (removes trailing zeros)
	return fmt.Sprintf("%g", node.Value)
}

// renderBinaryOp converts a BinaryOpNode to a string
func (r *Renderer) renderBinaryOp(node *BinaryOpNode, parentPrec int, isRightOperand bool) string {
	// Get the precedence of this operator
	opPrec := r.getPrecedence(node.Operator.Type)

	// Determine if we need parentheses
	needsParens := false
	if opPrec < parentPrec {
		// Lower precedence than parent, always needs parentheses
		needsParens = true
	} else if opPrec == parentPrec && isRightOperand {
		// Same precedence as parent, on the right side
		// Need parentheses for non-associative operators (-, /)
		if node.Operator.Type == MINUS || node.Operator.Type == DIVIDE {
			needsParens = true
		}
	}

	// Render left and right operands
	left := r.renderNode(node.Left, opPrec, false)
	right := r.renderNode(node.Right, opPrec, true)

	// Build the expression string
	var result string
	result = left + " " + node.Operator.Literal + " " + right

	if needsParens {
		result = "(" + result + ")"
	}

	return result
}

// renderGroup converts a GroupNode to a string
func (r *Renderer) renderGroup(node *GroupNode) string {
	inner := r.renderNode(node.Inner, LOWEST, false)
	return "(" + inner + ")"
}

// getPrecedence returns the precedence of an operator
func (r *Renderer) getPrecedence(tokenType TokenType) int {
	if prec, ok := precedences[tokenType]; ok {
		return prec
	}
	return LOWEST
}

// RenderLatex converts a LatexNode to a string
func RenderLatex(node LatexNode) string {
	renderer := NewRenderer()
	return renderer.Render(node)
}

// ExpressionToLatex converts an Expression tree directly to a LaTeX string
func ExpressionToLatex(expression expr.Expression) (string, error) {
	// First, export Expression to LaTeX AST
	ast, err := ExportToLatex(expression)
	if err != nil {
		return "", fmt.Errorf("failed to export expression: %w", err)
	}

	// Then, render the AST to a string
	latexString := RenderLatex(ast)

	return strings.TrimSpace(latexString), nil
}
