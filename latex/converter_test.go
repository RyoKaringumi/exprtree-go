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

	addExpr, ok := expression.(*expr.Add)
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

	subExpr, ok := expression.(*expr.Subtract)
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

	mulExpr, ok := expression.(*expr.Multiply)
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

	divExpr, ok := expression.(*expr.Divide)
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

func TestConvert_Power(t *testing.T) {
	// 2^3
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 2.0},
		Operator: Token{Type: CARET},
		Right:    &NumberNode{Value: 3.0},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	powExpr, ok := expression.(*expr.Power)
	if !ok {
		t.Fatalf("expected PowerExpression, got %T", expression)
	}

	// Verify it evaluates correctly: 2^3 = 8
	result, ok := powExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 8.0 {
		t.Errorf("expected result 8.0, got %f", numResult.Value)
	}
}

func TestConvert_SqrtBasic(t *testing.T) {
	// \sqrt{4}
	node := &CommandNode{
		Name:     "sqrt",
		Argument: &NumberNode{Value: 4.0},
		Optional: nil,
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	sqrtExpr, ok := expression.(*expr.Sqrt)
	if !ok {
		t.Fatalf("expected SqrtExpression, got %T", expression)
	}

	// Verify N is 2 (square root)
	if sqrtExpr.N != 2.0 {
		t.Errorf("expected N=2, got %f", sqrtExpr.N)
	}

	// Verify it evaluates correctly: sqrt(4) = 2
	result, ok := sqrtExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 2.0 {
		t.Errorf("expected result 2.0, got %f", numResult.Value)
	}
}

func TestConvert_SqrtWithOptional(t *testing.T) {
	// \sqrt[3]{8}
	node := &CommandNode{
		Name:     "sqrt",
		Argument: &NumberNode{Value: 8.0},
		Optional: &NumberNode{Value: 3.0},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	sqrtExpr, ok := expression.(*expr.Sqrt)
	if !ok {
		t.Fatalf("expected SqrtExpression, got %T", expression)
	}

	// Verify N is 3 (cube root)
	if sqrtExpr.N != 3.0 {
		t.Errorf("expected N=3, got %f", sqrtExpr.N)
	}

	// Verify it evaluates correctly: cbrt(8) = 2
	result, ok := sqrtExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok {
		t.Errorf("expected NumberValue result")
	}

	// Use tolerance for floating point comparison
	diff := numResult.Value - 2.0
	if diff < 0 {
		diff = -diff
	}
	if diff > 1e-10 {
		t.Errorf("expected result 2.0, got %f (diff: %e)", numResult.Value, diff)
	}
}

func TestConvert_PowerPrecedence(t *testing.T) {
	// 2 + 3^4 = 2 + 81 = 83
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 2.0},
		Operator: Token{Type: PLUS},
		Right: &BinaryOpNode{
			Left:     &NumberNode{Value: 3.0},
			Operator: Token{Type: CARET},
			Right:    &NumberNode{Value: 4.0},
		},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	// Verify it evaluates correctly
	result, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 83.0 {
		t.Errorf("expected result 83.0, got %f", numResult.Value)
	}
}

func TestConvert_PowerRightAssociative(t *testing.T) {
	// 2^3^2 = 2^(3^2) = 2^9 = 512
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 2.0},
		Operator: Token{Type: CARET},
		Right: &BinaryOpNode{
			Left:     &NumberNode{Value: 3.0},
			Operator: Token{Type: CARET},
			Right:    &NumberNode{Value: 2.0},
		},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	// Verify it evaluates correctly
	result, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 512.0 {
		t.Errorf("expected result 512.0, got %f", numResult.Value)
	}
}

func TestConvert_EqualBasic(t *testing.T) {
	// 2 = 2
	node := &EqualNode{
		Left:     &NumberNode{Value: 2.0},
		Operator: Token{Type: EQUAL},
		Right:    &NumberNode{Value: 2.0},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	equalExpr, ok := expression.(*expr.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	result, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := result.(*expr.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if !boolResult.Value {
		t.Errorf("expected true for 2 = 2")
	}
}

func TestConvert_EqualTrue(t *testing.T) {
	// 5 = 5
	node := &EqualNode{
		Left:     &NumberNode{Value: 5.0},
		Operator: Token{Type: EQUAL},
		Right:    &NumberNode{Value: 5.0},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	equalExpr, ok := expression.(*expr.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	result, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := result.(*expr.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if !boolResult.Value {
		t.Errorf("expected true for 5 = 5")
	}
}

func TestConvert_EqualFalse(t *testing.T) {
	// 2 = 3
	node := &EqualNode{
		Left:     &NumberNode{Value: 2.0},
		Operator: Token{Type: EQUAL},
		Right:    &NumberNode{Value: 3.0},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	equalExpr, ok := expression.(*expr.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	result, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := result.(*expr.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if boolResult.Value {
		t.Errorf("expected false for 2 = 3")
	}
}

func TestConvert_EqualFloatingPoint(t *testing.T) {
	// 0.1 + 0.2 = 0.3 (floating point tolerance test)
	node := &EqualNode{
		Left: &BinaryOpNode{
			Left:     &NumberNode{Value: 0.1},
			Operator: Token{Type: PLUS},
			Right:    &NumberNode{Value: 0.2},
		},
		Operator: Token{Type: EQUAL},
		Right:    &NumberNode{Value: 0.3},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	equalExpr, ok := expression.(*expr.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly with floating point tolerance
	result, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := result.(*expr.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if !boolResult.Value {
		t.Errorf("expected true for 0.1+0.2=0.3 (with floating point tolerance)")
	}
}

func TestConvert_EqualComplex(t *testing.T) {
	// 2 + 3 = 1 + 4
	node := &EqualNode{
		Left: &BinaryOpNode{
			Left:     &NumberNode{Value: 2.0},
			Operator: Token{Type: PLUS},
			Right:    &NumberNode{Value: 3.0},
		},
		Operator: Token{Type: EQUAL},
		Right: &BinaryOpNode{
			Left:     &NumberNode{Value: 1.0},
			Operator: Token{Type: PLUS},
			Right:    &NumberNode{Value: 4.0},
		},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	equalExpr, ok := expression.(*expr.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly: (2+3) = (1+4) -> 5 = 5 -> true
	result, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := result.(*expr.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if !boolResult.Value {
		t.Errorf("expected true for (2+3)=(1+4)")
	}
}

func TestConvert_EqualWithGroups(t *testing.T) {
	// (2 + 3) = (1 + 4)
	node := &EqualNode{
		Left: &GroupNode{
			Inner: &BinaryOpNode{
				Left:     &NumberNode{Value: 2.0},
				Operator: Token{Type: PLUS},
				Right:    &NumberNode{Value: 3.0},
			},
		},
		Operator: Token{Type: EQUAL},
		Right: &GroupNode{
			Inner: &BinaryOpNode{
				Left:     &NumberNode{Value: 1.0},
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

	equalExpr, ok := expression.(*expr.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	result, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := result.(*expr.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if !boolResult.Value {
		t.Errorf("expected true for (2+3)=(1+4)")
	}
}

func TestConvert_EqualNested(t *testing.T) {
	// (2 = 2) = (3 = 3) -> true = true -> true
	node := &EqualNode{
		Left: &EqualNode{
			Left:     &NumberNode{Value: 2.0},
			Operator: Token{Type: EQUAL},
			Right:    &NumberNode{Value: 2.0},
		},
		Operator: Token{Type: EQUAL},
		Right: &EqualNode{
			Left:     &NumberNode{Value: 3.0},
			Operator: Token{Type: EQUAL},
			Right:    &NumberNode{Value: 3.0},
		},
	}
	converter := NewConverter()

	expression, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	equalExpr, ok := expression.(*expr.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	result, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := result.(*expr.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if !boolResult.Value {
		t.Errorf("expected true for (2=2)=(3=3)")
	}
}
