package latex

import (
	"exprtree/expr"
	"testing"
)

func TestConvert_Number(t *testing.T) {
	node := &NumberNode{Value: 42.0}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	constant, ok := expression.(*expr.Constant)
	if !ok {
		t.Fatalf("expected Constant, got %T", expression)
	}

	if constant.Value.Value != 42.0 {
		t.Errorf("expected value 42.0, got %f", constant.Value.Value)
	}
}

func TestConvert_Addition(t *testing.T) {
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 2.0},
		Operator: Token{Type: PLUS},
		Right:    &NumberNode{Value: 3.0},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	addExpr, ok := expression.(*expr.AddExpression)
	if !ok {
		t.Fatalf("expected AddExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	result, ok := addExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 5.0 {
		t.Errorf("expected result 5.0, got %f", numResult.Value)
	}
}

func TestConvert_Subtraction(t *testing.T) {
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 10.0},
		Operator: Token{Type: MINUS},
		Right:    &NumberNode{Value: 3.0},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	subExpr, ok := expression.(*expr.SubtractExpression)
	if !ok {
		t.Fatalf("expected SubtractExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	result, ok := subExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 7.0 {
		t.Errorf("expected result 7.0, got %f", numResult.Value)
	}
}

func TestConvert_Multiplication(t *testing.T) {
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 6.0},
		Operator: Token{Type: MULTIPLY},
		Right:    &NumberNode{Value: 7.0},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	mulExpr, ok := expression.(*expr.MultiplyExpression)
	if !ok {
		t.Fatalf("expected MultiplyExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	result, ok := mulExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 42.0 {
		t.Errorf("expected result 42.0, got %f", numResult.Value)
	}
}

func TestConvert_Division(t *testing.T) {
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 15.0},
		Operator: Token{Type: DIVIDE},
		Right:    &NumberNode{Value: 3.0},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	divExpr, ok := expression.(*expr.DivideExpression)
	if !ok {
		t.Fatalf("expected DivideExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	result, ok := divExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 5.0 {
		t.Errorf("expected result 5.0, got %f", numResult.Value)
	}
}

func TestConvert_Precedence(t *testing.T) {
	// 2 + (3 * 4)
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 2.0},
		Operator: Token{Type: PLUS},
		Right: &BinaryOpNode{
			Left:     &NumberNode{Value: 3.0},
			Operator: Token{Type: MULTIPLY},
			Right:    &NumberNode{Value: 4.0},
		},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	// Verify it evaluates correctly: 2 + 12 = 14
	result, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 14.0 {
		t.Errorf("expected result 14.0, got %f", numResult.Value)
	}
}

func TestConvert_Group(t *testing.T) {
	// (2 + 3) * 4
	node := &BinaryOpNode{
		Left: &GroupNode{
			Inner: &BinaryOpNode{
				Left:     &NumberNode{Value: 2.0},
				Operator: Token{Type: PLUS},
				Right:    &NumberNode{Value: 3.0},
			},
		},
		Operator: Token{Type: MULTIPLY},
		Right:    &NumberNode{Value: 4.0},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	// Verify it evaluates correctly: 5 * 4 = 20
	result, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 20.0 {
		t.Errorf("expected result 20.0, got %f", numResult.Value)
	}
}

func TestConvert_ComplexTree(t *testing.T) {
	// (1 + 2) * (3 + 4)
	node := &BinaryOpNode{
		Left: &GroupNode{
			Inner: &BinaryOpNode{
				Left:     &NumberNode{Value: 1.0},
				Operator: Token{Type: PLUS},
				Right:    &NumberNode{Value: 2.0},
			},
		},
		Operator: Token{Type: MULTIPLY},
		Right: &GroupNode{
			Inner: &BinaryOpNode{
				Left:     &NumberNode{Value: 3.0},
				Operator: Token{Type: PLUS},
				Right:    &NumberNode{Value: 4.0},
			},
		},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	// Verify it evaluates correctly: 3 * 7 = 21
	result, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 21.0 {
		t.Errorf("expected result 21.0, got %f", numResult.Value)
	}
}

func TestConvert_NilNode(t *testing.T) {
	converter := NewConverter()

	_, err := converter.Convert(nil)
	if err == nil {
		t.Errorf("expected error for nil node")
	}
}
