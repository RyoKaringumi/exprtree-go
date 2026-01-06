package latex

import (
	"exprtree/expr"
	"testing"
)

func TestExportConstant(t *testing.T) {
	// Create a constant expression
	constant := &expr.Constant{
		Value: expr.NumberValue{Value: 42.5},
	}

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
	left := &expr.Constant{Value: expr.NumberValue{Value: 2}}
	right := &expr.Constant{Value: expr.NumberValue{Value: 3}}
	addExpr := expr.NewAddExpression(left, right)

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
	left := expr.NewAddExpression(
		&expr.Constant{Value: expr.NumberValue{Value: 2}},
		&expr.Constant{Value: expr.NumberValue{Value: 3}},
	)
	right := &expr.Constant{Value: expr.NumberValue{Value: 4}}
	multExpr := expr.NewMultiplyExpression(left, right)

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
		expr     expr.Expression
		expected TokenType
	}{
		{
			name:     "Addition",
			expr:     expr.NewAddExpression(&expr.Constant{Value: expr.NumberValue{Value: 1}}, &expr.Constant{Value: expr.NumberValue{Value: 2}}),
			expected: PLUS,
		},
		{
			name:     "Subtraction",
			expr:     expr.NewSubtractExpression(&expr.Constant{Value: expr.NumberValue{Value: 5}}, &expr.Constant{Value: expr.NumberValue{Value: 3}}),
			expected: MINUS,
		},
		{
			name:     "Multiplication",
			expr:     expr.NewMultiplyExpression(&expr.Constant{Value: expr.NumberValue{Value: 2}}, &expr.Constant{Value: expr.NumberValue{Value: 3}}),
			expected: MULTIPLY,
		},
		{
			name:     "Division",
			expr:     expr.NewDivideExpression(&expr.Constant{Value: expr.NumberValue{Value: 10}}, &expr.Constant{Value: expr.NumberValue{Value: 2}}),
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
