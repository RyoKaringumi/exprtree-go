package latex

import (
	"exprtree/expr"
	"exprtree/value"
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
func (e *Exporter) Export(expression expr.Expr) (LatexNode, error) {
	if expression == nil {
		return nil, fmt.Errorf("cannot export nil expression")
	}

	switch exp := expression.(type) {
	case *expr.Constant:
		return e.exportConstant(exp), nil
	case *expr.Add:
		return e.exportBinaryOp(exp, PLUS, "+")
	case *expr.Sub:
		return e.exportBinaryOp(exp, MINUS, "-")
	case *expr.Mul:
		return e.exportBinaryOp(exp, MULTIPLY, "*")
	case *expr.Div:
		return e.exportBinaryOp(exp, DIVIDE, "/")
	case *expr.Variable:
		return e.exportVariable(exp), nil
	default:
		return nil, fmt.Errorf("unknown expression type: %T", expression)
	}
}

// exportConstant converts a Constant to a NumberNode
func (e *Exporter) exportConstant(constant *expr.Constant) LatexNode {
	// Float64() to get the float64 value
	realValue, ok := constant.Value().(*value.RealValue)
	if !ok {
		e.errors = append(e.errors, fmt.Sprintf("unsupported constant value type: %T", constant.Value()))
		return &NumberNode{
			Value: 0,
			Token: Token{
				Type:    NUMBER,
				Literal: "0",
				Value:   0,
			},
		}
	}

	return &NumberNode{
		Value: realValue.Float64(),
		Token: Token{
			Type:    NUMBER,
			Literal: fmt.Sprintf("%g", realValue.Float64()),
			Value:   realValue.Float64(),
		},
	}
}

// exportVariable converts a Variable to a VariableNode
func (e *Exporter) exportVariable(variable *expr.Variable) LatexNode {
	name := variable.Name()
	return &VariableNode{
		Name: name,
		Token: Token{
			Type:    VARIABLE,
			Literal: name,
		},
	}
}

// exportBinaryOp converts a binary expression to a BinaryOpNode
func (e *Exporter) exportBinaryOp(binaryExpr expr.Expr, opType TokenType, opLiteral string) (LatexNode, error) {
	binary, ok := binaryExpr.(expr.Binary)
	if !ok {
		return nil, fmt.Errorf("expression is not binary: %T", binaryExpr)
	}
	left, err := e.Export(binary.Left())
	if err != nil {
		return nil, fmt.Errorf("failed to export left operand: %w", err)
	}
	right, err := e.Export(binary.Right())
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
func ExportToLatex(expression expr.Expr) (LatexNode, error) {
	exporter := NewExporter()
	return exporter.Export(expression)
}
