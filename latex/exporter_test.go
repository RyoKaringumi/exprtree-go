package latex

import (
	"exprtree/expr"
	"exprtree/value"
	"testing"
)

func TestExportConstant(t *testing.T) {
	// Create a constant expression
	constant := expr.NewConstant(value.NewRealValue(42.5))

	// Export to LaTeX AST
	exporter := NewExporter()
	node, err := exporter.Export(constant)

	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify it's a NumberNode
	numNode, ok := node.(*NumberNode)
	if !ok {
		t.Fatalf("Expected NumberNode, got %T", node)
	}

	if numNode.Value != 42.5 {
		t.Errorf("Expected value 42.5, got %f", numNode.Value)
	}
}

func TestExportAddExpression(t *testing.T) {
	// Create: 2 + 3
	left := expr.NewConstant(value.NewRealValue(2))
	right := expr.NewConstant(value.NewRealValue(3))
	addExpr := expr.NewAdd(left, right)

	// Export to LaTeX AST
	exporter := NewExporter()
	node, err := exporter.Export(addExpr)

	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify it's a BinaryOpNode
	binNode, ok := node.(*BinaryOpNode)
	if !ok {
		t.Fatalf("Expected BinaryOpNode, got %T", node)
	}

	if binNode.Operator.Type != PLUS {
		t.Errorf("Expected PLUS operator, got %v", binNode.Operator.Type)
	}

	// Verify left operand
	leftNode, ok := binNode.Left.(*NumberNode)
	if !ok || leftNode.Value != 2 {
		t.Errorf("Expected left operand to be NumberNode with value 2")
	}

	// Verify right operand
	rightNode, ok := binNode.Right.(*NumberNode)
	if !ok || rightNode.Value != 3 {
		t.Errorf("Expected right operand to be NumberNode with value 3")
	}
}

func TestExportComplexExpression(t *testing.T) {
	// Create: (2 + 3) * 4
	left := expr.NewAdd(
		expr.NewConstant(value.NewRealValue(2)),
		expr.NewConstant(value.NewRealValue(3)),
	)
	right := expr.NewConstant(value.NewRealValue(4))
	multExpr := expr.NewMul(left, right)

	// Export to LaTeX AST
	exporter := NewExporter()
	node, err := exporter.Export(multExpr)

	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}

	// Verify top-level is multiply
	binNode, ok := node.(*BinaryOpNode)
	if !ok || binNode.Operator.Type != MULTIPLY {
		t.Fatalf("Expected BinaryOpNode with MULTIPLY operator")
	}

	// Verify left is add operation
	leftBinNode, ok := binNode.Left.(*BinaryOpNode)
	if !ok || leftBinNode.Operator.Type != PLUS {
		t.Errorf("Expected left operand to be BinaryOpNode with PLUS operator")
	}
}

func TestExportAllOperators(t *testing.T) {
	tests := []struct {
		name     string
		expr     expr.Expr
		expected TokenType
	}{
		{
			name:     "Addition",
			expr:     expr.NewAdd(expr.NewConstant(value.NewRealValue(1)), expr.NewConstant(value.NewRealValue(2))),
			expected: PLUS,
		},
		{
			name:     "Subtraction",
			expr:     expr.NewSub(expr.NewConstant(value.NewRealValue(5)), expr.NewConstant(value.NewRealValue(3))),
			expected: MINUS,
		},
		{
			name:     "Multiplication",
			expr:     expr.NewMul(expr.NewConstant(value.NewRealValue(2)), expr.NewConstant(value.NewRealValue(3))),
			expected: MULTIPLY,
		},
		{
			name:     "Division",
			expr:     expr.NewDiv(expr.NewConstant(value.NewRealValue(10)), expr.NewConstant(value.NewRealValue(2))),
			expected: DIVIDE,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exporter := NewExporter()
			node, err := exporter.Export(tt.expr)

			if err != nil {
				t.Fatalf("Export failed: %v", err)
			}

			binNode, ok := node.(*BinaryOpNode)
			if !ok {
				t.Fatalf("Expected BinaryOpNode, got %T", node)
			}

			if binNode.Operator.Type != tt.expected {
				t.Errorf("Expected operator %v, got %v", tt.expected, binNode.Operator.Type)
			}
		})
	}
}
