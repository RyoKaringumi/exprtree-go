package latex

import (
	"exprtree/expr"
	"testing"
)

func TestRenderNumber(t *testing.T) {
	node := &NumberNode{
		Value: 42.5,
	}

	renderer := NewRenderer()
	result := renderer.Render(node)

	expected := "42.5"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestRenderSimpleBinaryOp(t *testing.T) {
	tests := []struct {
		name     string
		left     float64
		op       TokenType
		opLit    string
		right    float64
		expected string
	}{
		{"Addition", 2, PLUS, "+", 3, "2 + 3"},
		{"Subtraction", 5, MINUS, "-", 3, "5 - 3"},
		{"Multiplication", 2, MULTIPLY, "*", 3, "2 * 3"},
		{"Division", 10, DIVIDE, "/", 2, "10 / 2"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &BinaryOpNode{
				Left:  &NumberNode{Value: tt.left},
				Operator: Token{
					Type:    tt.op,
					Literal: tt.opLit,
				},
				Right: &NumberNode{Value: tt.right},
			}

			renderer := NewRenderer()
			result := renderer.Render(node)

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestRenderPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		node     LatexNode
		expected string
	}{
		{
			name: "Multiplication before addition (left)",
			node: &BinaryOpNode{
				Left: &BinaryOpNode{
					Left:     &NumberNode{Value: 2},
					Operator: Token{Type: MULTIPLY, Literal: "*"},
					Right:    &NumberNode{Value: 3},
				},
				Operator: Token{Type: PLUS, Literal: "+"},
				Right:    &NumberNode{Value: 4},
			},
			expected: "2 * 3 + 4",
		},
		{
			name: "Addition before multiplication (needs parentheses)",
			node: &BinaryOpNode{
				Left: &BinaryOpNode{
					Left:     &NumberNode{Value: 2},
					Operator: Token{Type: PLUS, Literal: "+"},
					Right:    &NumberNode{Value: 3},
				},
				Operator: Token{Type: MULTIPLY, Literal: "*"},
				Right:    &NumberNode{Value: 4},
			},
			expected: "(2 + 3) * 4",
		},
		{
			name: "Multiplication before addition (right)",
			node: &BinaryOpNode{
				Left:     &NumberNode{Value: 2},
				Operator: Token{Type: PLUS, Literal: "+"},
				Right: &BinaryOpNode{
					Left:     &NumberNode{Value: 3},
					Operator: Token{Type: MULTIPLY, Literal: "*"},
					Right:    &NumberNode{Value: 4},
				},
			},
			expected: "2 + 3 * 4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer()
			result := renderer.Render(tt.node)

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestRenderAssociativity(t *testing.T) {
	tests := []struct {
		name     string
		node     LatexNode
		expected string
	}{
		{
			name: "Left associative subtraction",
			node: &BinaryOpNode{
				Left: &BinaryOpNode{
					Left:     &NumberNode{Value: 10},
					Operator: Token{Type: MINUS, Literal: "-"},
					Right:    &NumberNode{Value: 3},
				},
				Operator: Token{Type: MINUS, Literal: "-"},
				Right:    &NumberNode{Value: 2},
			},
			expected: "10 - 3 - 2",
		},
		{
			name: "Right associative subtraction (needs parentheses)",
			node: &BinaryOpNode{
				Left:     &NumberNode{Value: 10},
				Operator: Token{Type: MINUS, Literal: "-"},
				Right: &BinaryOpNode{
					Left:     &NumberNode{Value: 3},
					Operator: Token{Type: MINUS, Literal: "-"},
					Right:    &NumberNode{Value: 2},
				},
			},
			expected: "10 - (3 - 2)",
		},
		{
			name: "Left associative division",
			node: &BinaryOpNode{
				Left: &BinaryOpNode{
					Left:     &NumberNode{Value: 12},
					Operator: Token{Type: DIVIDE, Literal: "/"},
					Right:    &NumberNode{Value: 4},
				},
				Operator: Token{Type: DIVIDE, Literal: "/"},
				Right:    &NumberNode{Value: 2},
			},
			expected: "12 / 4 / 2",
		},
		{
			name: "Right associative division (needs parentheses)",
			node: &BinaryOpNode{
				Left:     &NumberNode{Value: 12},
				Operator: Token{Type: DIVIDE, Literal: "/"},
				Right: &BinaryOpNode{
					Left:     &NumberNode{Value: 4},
					Operator: Token{Type: DIVIDE, Literal: "/"},
					Right:    &NumberNode{Value: 2},
				},
			},
			expected: "12 / (4 / 2)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			renderer := NewRenderer()
			result := renderer.Render(tt.node)

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestExpressionToLatex(t *testing.T) {
	tests := []struct {
		name     string
		expr     expr.Expression
		expected string
	}{
		{
			name:     "Simple constant",
			expr:     &expr.Constant{Value: expr.NumberValue{Value: 42}},
			expected: "42",
		},
		{
			name:     "Simple addition",
			expr:     expr.NewAddExpression(&expr.Constant{Value: expr.NumberValue{Value: 2}}, &expr.Constant{Value: expr.NumberValue{Value: 3}}),
			expected: "2 + 3",
		},
		{
			name: "Complex expression with precedence",
			expr: expr.NewMultiplyExpression(
				expr.NewAddExpression(
					&expr.Constant{Value: expr.NumberValue{Value: 2}},
					&expr.Constant{Value: expr.NumberValue{Value: 3}},
				),
				&expr.Constant{Value: expr.NumberValue{Value: 4}},
			),
			expected: "(2 + 3) * 4",
		},
		{
			name: "Nested operations",
			expr: expr.NewAddExpression(
				&expr.Constant{Value: expr.NumberValue{Value: 2}},
				expr.NewMultiplyExpression(
					&expr.Constant{Value: expr.NumberValue{Value: 3}},
					&expr.Constant{Value: expr.NumberValue{Value: 4}},
				),
			),
			expected: "2 + 3 * 4",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ExpressionToLatex(tt.expr)

			if err != nil {
				t.Fatalf("ExpressionToLatex failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestRoundTrip(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"Simple addition", "2 + 3"},
		{"Multiplication", "2 * 3"},
		{"Precedence", "2 + 3 * 4"},
		{"Parentheses", "(2 + 3) * 4"},
		{"Left associativity", "10 - 3 - 2"},
		{"Right associativity", "10 - (3 - 2)"},
		{"Complex expression", "(1 + 2) * (3 + 4)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse the input
			expression, err := ParseLatex(tt.input)
			if err != nil {
				t.Fatalf("ParseLatex failed: %v", err)
			}

			// Convert back to string
			result, err := ExpressionToLatex(expression)
			if err != nil {
				t.Fatalf("ExpressionToLatex failed: %v", err)
			}

			// Parse the result again
			expression2, err := ParseLatex(result)
			if err != nil {
				t.Fatalf("Second ParseLatex failed: %v", err)
			}

			// Evaluate both and compare
			val1, ok1 := expression.Eval()
			val2, ok2 := expression2.Eval()

			if !ok1 || !ok2 {
				t.Fatalf("Evaluation failed")
			}

			num1, ok1 := val1.(*expr.NumberValue)
			num2, ok2 := val2.(*expr.NumberValue)

			if !ok1 || !ok2 {
				t.Fatalf("Result is not a number")
			}

			if num1.Value != num2.Value {
				t.Errorf("Values differ: %f != %f (input: %s, output: %s)", num1.Value, num2.Value, tt.input, result)
			}
		})
	}
}
