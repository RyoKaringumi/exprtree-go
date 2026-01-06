package latex

import (
	"exprtree/expr"
	"fmt"
)

// Exporter converts Expression tree to LaTeX AST
type Exporter struct {
	errors []string
}

// NewExporter creates a new Exporter instance
func NewExporter() *Exporter {
	return &Exporter{
		errors: []string{},
	}
}

// Export converts an Expression to a LatexNode
func (e *Exporter) Export(expression expr.Expression) (LatexNode, error) {
	if expression == nil {
		return nil, fmt.Errorf("cannot export nil expression")
	}

	switch exp := expression.(type) {
	case *expr.Constant:
		return e.exportConstant(exp), nil
	case *expr.AddExpression:
		return e.exportBinaryOp(exp, PLUS, "+")
	case *expr.SubtractExpression:
		return e.exportBinaryOp(exp, MINUS, "-")
	case *expr.MultiplyExpression:
		return e.exportBinaryOp(exp, MULTIPLY, "*")
	case *expr.DivideExpression:
		return e.exportBinaryOp(exp, DIVIDE, "/")
	case *expr.Variable:
		return e.exportVariable(exp), nil
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expression)
	}
}

// exportConstant converts a Constant to a NumberNode
func (e *Exporter) exportConstant(constant *expr.Constant) LatexNode {
	return &NumberNode{
		Value: constant.Value.Value,
		Token: Token{
			Type:    NUMBER,
			Literal: fmt.Sprintf("%g", constant.Value.Value),
			Value:   constant.Value.Value,
		},
	}
}

// exportVariable converts a Variable to a VariableNode
func (e *Exporter) exportVariable(variable *expr.Variable) LatexNode {
	return &VariableNode{
		Name: variable.Name,
		Token: Token{
			Type:    VARIABLE,
			Literal: variable.Name,
		},
	}
}

// exportBinaryOp converts a binary expression to a BinaryOpNode
func (e *Exporter) exportBinaryOp(binaryExpr expr.Expression, opType TokenType, opLiteral string) (LatexNode, error) {
	children := binaryExpr.Children()
	if len(children) != 2 {
		return nil, fmt.Errorf("binary expression must have exactly 2 children")
	}

	left, err := e.Export(children[0])
	if err != nil {
		return nil, fmt.Errorf("failed to export left operand: %w", err)
	}

	right, err := e.Export(children[1])
	if err != nil {
		return nil, fmt.Errorf("failed to export right operand: %w", err)
	}

	return &BinaryOpNode{
		Left: left,
		Operator: Token{
			Type:    opType,
			Literal: opLiteral,
		},
		Right: right,
	}, nil
}

// Errors returns the list of export errors
func (e *Exporter) Errors() []string {
	return e.errors
}

// ExportToLatex converts an Expression tree to a LaTeX AST
func ExportToLatex(expression expr.Expression) (LatexNode, error) {
	exporter := NewExporter()
	return exporter.Export(expression)
}
