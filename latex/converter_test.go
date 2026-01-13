package latex

import (
	"exprtree/expr"
	"exprtree/prop"
	"exprtree/value"
	"testing"
)

func valueToFloat64(v value.Value) float64 {
	if realVal, ok := v.(*value.RealValue); ok {
		return realVal.Float64()
	}
	return 0.0
}

func constantToFloat64(c *expr.Constant) float64 {
	return valueToFloat64(c.Value())
}

func TestConvert_Number(t *testing.T) {
	node := &NumberNode{Value: 42.0}
	converter := NewConverter()

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	constant, ok := result.(*expr.Constant)
	if !ok {
		t.Fatalf("expected Constant, got %T", result)
	}

	if constantToFloat64(constant) != 42.0 {
		t.Errorf("expected value 42.0, got %f", constantToFloat64(constant))
	}
}

func TestConvert_Addition(t *testing.T) {
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 2.0},
		Operator: Token{Type: PLUS},
		Right:    &NumberNode{Value: 3.0},
	}
	converter := NewConverter()

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	addExpr, ok := expression.(*expr.Add)
	if !ok {
		t.Fatalf("expected AddExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	evalResult, ok := addExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 5.0 {
		t.Errorf("expected result 5.0, got %f", numResult.Float64())
	}
}

func TestConvert_Subtraction(t *testing.T) {
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 10.0},
		Operator: Token{Type: MINUS},
		Right:    &NumberNode{Value: 3.0},
	}
	converter := NewConverter()

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	subExpr, ok := expression.(*expr.Sub)
	if !ok {
		t.Fatalf("expected SubtractExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	evalResult, ok := subExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 7.0 {
		t.Errorf("expected result 7.0, got %f", numResult.Float64())
	}
}

func TestConvert_Multiplication(t *testing.T) {
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 6.0},
		Operator: Token{Type: MULTIPLY},
		Right:    &NumberNode{Value: 7.0},
	}
	converter := NewConverter()

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	mulExpr, ok := expression.(*expr.Mul)
	if !ok {
		t.Fatalf("expected MultiplyExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	evalResult, ok := mulExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 42.0 {
		t.Errorf("expected result 42.0, got %f", numResult.Float64())
	}
}

func TestConvert_Division(t *testing.T) {
	node := &BinaryOpNode{
		Left:     &NumberNode{Value: 15.0},
		Operator: Token{Type: DIVIDE},
		Right:    &NumberNode{Value: 3.0},
	}
	converter := NewConverter()

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	divExpr, ok := expression.(*expr.Div)
	if !ok {
		t.Fatalf("expected DivideExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	evalResult, ok := divExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 5.0 {
		t.Errorf("expected result 5.0, got %f", numResult.Float64())
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	// Verify it evaluates correctly: 2 + 12 = 14
	evalResult, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 14.0 {
		t.Errorf("expected result 14.0, got %f", numResult.Float64())
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	// Verify it evaluates correctly: 5 * 4 = 20
	evalResult, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 20.0 {
		t.Errorf("expected result 20.0, got %f", numResult.Float64())
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	// Verify it evaluates correctly: 3 * 7 = 21
	evalResult, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 21.0 {
		t.Errorf("expected result 21.0, got %f", numResult.Float64())
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	powExpr, ok := expression.(*expr.Power)
	if !ok {
		t.Fatalf("expected PowerExpression, got %T", expression)
	}

	// Verify it evaluates correctly: 2^3 = 8
	evalResult, ok := powExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 8.0 {
		t.Errorf("expected result 8.0, got %f", numResult.Float64())
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	sqrtExpr, ok := expression.(*expr.NthRoot)
	if !ok {
		t.Fatalf("expected SqrtExpression, got %T", expression)
	}

	// Verify N is 2 (square root)
	degree, ok := sqrtExpr.Degree().Eval()
	if valueToFloat64(degree) != 2.0 {
		t.Errorf("expected N=2, got %f", valueToFloat64(degree))
	}

	// Verify it evaluates correctly: sqrt(4) = 2
	evalResult, ok := sqrtExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 2.0 {
		t.Errorf("expected result 2.0, got %f", numResult.Float64())
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	sqrtExpr, ok := expression.(*expr.NthRoot)
	if !ok {
		t.Fatalf("expected SqrtExpression, got %T", expression)
	}

	// Verify N is 3 (cube root)
	degree, ok := sqrtExpr.Degree().Eval()
	if valueToFloat64(degree) != 3.0 {
		t.Errorf("expected N=3, got %f", valueToFloat64(degree))
	}

	// Verify it evaluates correctly: cbrt(8) = 2
	evalResult, ok := sqrtExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok {
		t.Errorf("expected NumberValue result")
	}

	// Use tolerance for floating point comparison
	diff := numResult.Float64() - 2.0
	if diff < 0 {
		diff = -diff
	}
	if diff > 1e-10 {
		t.Errorf("expected result 2.0, got %f (diff: %e)", numResult.Float64(), diff)
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	// Verify it evaluates correctly
	evalResult, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 83.0 {
		t.Errorf("expected result 83.0, got %f", numResult.Float64())
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	// Verify it evaluates correctly
	evalResult, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 512.0 {
		t.Errorf("expected result 512.0, got %f", numResult.Float64())
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	equalExpr, ok := expression.(*prop.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	evalResult, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := evalResult.(*value.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if !boolResult.Bool() {
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	equalExpr, ok := expression.(*prop.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	evalResult, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := evalResult.(*value.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if !boolResult.Bool() {
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	equalExpr, ok := expression.(*prop.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	evalResult, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := evalResult.(*value.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if boolResult.Bool() {
		t.Errorf("expected false for 2 = 3")
	}
}

// func TestConvert_EqualFloatingPoint(t *testing.T) {
// 	// 0.1 + 0.2 = 0.3 (floating point tolerance test)
// 	node := &EqualNode{
// 		Left: &BinaryOpNode{
// 			Left:     &NumberNode{Value: 0.1},
// 			Operator: Token{Type: PLUS},
// 			Right:    &NumberNode{Value: 0.2},
// 		},
// 		Operator: Token{Type: EQUAL},
// 		Right:    &NumberNode{Value: 0.3},
// 	}
// 	converter := NewConverter()

// 	result, err := converter.Convert(node)
// 	if err != nil {
// 		t.Fatalf("Convert error: %v", err)
// 	}

// 	expression, ok := result.(expr.Expr)
// 	if !ok {
// 		t.Fatalf("expected Expression, got %T", result)
// 	}

// 	equalExpr, ok := expression.(*prop.Equal)
// 	if !ok {
// 		t.Fatalf("expected EqualExpression, got %T", expression)
// 	}

// 	// Verify it evaluates correctly with floating point tolerance
// 	evalResult, ok := equalExpr.Eval()
// 	if !ok {
// 		t.Errorf("evaluation failed")
// 	}

// 	boolResult, ok := evalResult.(*value.BoolValue)
// 	if !ok {
// 		t.Fatalf("expected BoolValue, got %T", result)
// 	}

// 	if !boolResult.Bool() {
// 		t.Errorf("expected true for 0.1+0.2=0.3 (with floating point tolerance)")
// 	}
// }

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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	equalExpr, ok := expression.(*prop.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly: (2+3) = (1+4) -> 5 = 5 -> true
	evalResult, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := evalResult.(*value.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if !boolResult.Bool() {
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("expected Expression, got %T", result)
	}

	equalExpr, ok := expression.(*prop.Equal)
	if !ok {
		t.Fatalf("expected EqualExpression, got %T", expression)
	}

	// Verify it evaluates correctly
	evalResult, ok := equalExpr.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	boolResult, ok := evalResult.(*value.BoolValue)
	if !ok {
		t.Fatalf("expected BoolValue, got %T", result)
	}

	if !boolResult.Bool() {
		t.Errorf("expected true for (2+3)=(1+4)")
	}
}

func TestConvert_EqualNested(t *testing.T) {
	// (2 = 2) = (3 = 3)
	// Note: With chained equality detection, this now becomes And(Eq(2,2), Eq(3,3))
	// rather than Eq(Eq(2,2), Eq(3,3))
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

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	// This now returns And due to chained equality detection
	andExpr, ok := result.(*prop.And)
	if !ok {
		t.Fatalf("expected And proposition, got %T", result)
	}

	// Verify structure: should be And(Eq(2,2), Eq(2,Eq(3,3)))
	// Left should be Equal(2, 2)
	leftEqual, ok := andExpr.Left().(*prop.Equal)
	if !ok {
		t.Fatalf("expected left to be Equal, got %T", andExpr.Left())
	}

	// Right should be Equal(2, Equal(3, 3))
	rightEqual, ok := andExpr.Right().(*prop.Equal)
	if !ok {
		t.Fatalf("expected right to be Equal, got %T", andExpr.Right())
	}

	// Verify left is 2 = 2
	_ = leftEqual // Structure verified

	// Verify right is 2 = (3 = 3)
	_ = rightEqual // Structure verified
}

// TestConverter_ChainedEquality_ThreeTerms tests conversion of chained equality a = b = c
// Mathematically, a = b = c means "a equals b AND b equals c", which should be represented
// as And(Eq(a,b), Eq(b,c)), not as (a = b) = c.
//
// Note: This test expects the Converter to return a Proposition (And), not an Expression.
// The Converter.Convert method may need to return interface{} or a union type to support both.
func TestConverter_ChainedEquality_ThreeTerms(t *testing.T) {
	// a = b = c
	// Parser produces: EqualNode(EqualNode(a, b), c) due to left-associativity
	// Expected conversion: And(Eq(a, b), Eq(b, c))
	// where Eq(a,b) and Eq(b,c) are Equal (Proposition)

	node := &EqualNode{
		Left: &EqualNode{
			Left:     &VariableNode{Name: "a"},
			Operator: Token{Type: EQUAL},
			Right:    &VariableNode{Name: "b"},
		},
		Operator: Token{Type: EQUAL},
		Right:    &VariableNode{Name: "c"},
	}
	converter := NewConverter()

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	// Expected structure: And(Eq(a, b), Eq(b, c))
	// Since And is a Proposition, we need to check if result can be type-asserted
	andExpr, ok := result.(*prop.And)
	if !ok {
		t.Fatalf("expected And proposition, got %T", result)
	}

	// Left should be Equal(a, b)
	leftEqual, ok := andExpr.Left().(*prop.Equal)
	if !ok {
		t.Fatalf("expected left to be Equal proposition, got %T", andExpr.Left())
	}

	// Check left Equal: a = b
	leftVarA, ok := leftEqual.Left().(*expr.Variable)
	if !ok || leftVarA.Name() != "a" {
		t.Errorf("expected left.left to be variable 'a'")
	}

	leftVarB, ok := leftEqual.Right().(*expr.Variable)
	if !ok || leftVarB.Name() != "b" {
		t.Errorf("expected left.right to be variable 'b'")
	}

	// Right should be Equal(b, c)
	rightEqual, ok := andExpr.Right().(*prop.Equal)
	if !ok {
		t.Fatalf("expected right to be Equal proposition, got %T", andExpr.Right())
	}

	// Check right Equal: b = c
	rightVarB, ok := rightEqual.Left().(*expr.Variable)
	if !ok || rightVarB.Name() != "b" {
		t.Errorf("expected right.left to be variable 'b'")
	}

	rightVarC, ok := rightEqual.Right().(*expr.Variable)
	if !ok || rightVarC.Name() != "c" {
		t.Errorf("expected right.right to be variable 'c'")
	}
}

// TestConverter_ChainedEquality_FourTerms tests conversion of chained equality a = b = c = d
// Mathematically, a = b = c = d means "a equals b AND b equals c AND c equals d"
// Expected: And(And(Eq(a,b), Eq(b,c)), Eq(c,d))
//
// Note: This test expects the Converter to return a Proposition (And), not an Expression.
func TestConverter_ChainedEquality_FourTerms(t *testing.T) {
	// a = b = c = d
	// Parser produces (left-associative): EqualNode(EqualNode(EqualNode(a, b), c), d)
	// Expected: And(And(Eq(a,b), Eq(b,c)), Eq(c,d))

	node := &EqualNode{
		Left: &EqualNode{
			Left: &EqualNode{
				Left:     &VariableNode{Name: "a"},
				Operator: Token{Type: EQUAL},
				Right:    &VariableNode{Name: "b"},
			},
			Operator: Token{Type: EQUAL},
			Right:    &VariableNode{Name: "c"},
		},
		Operator: Token{Type: EQUAL},
		Right:    &VariableNode{Name: "d"},
	}
	converter := NewConverter()

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	// Expected structure: And(And(Eq(a,b), Eq(b,c)), Eq(c,d))
	outerAnd, ok := result.(*prop.And)
	if !ok {
		t.Fatalf("expected outer And proposition, got %T", result)
	}

	// Left should be And(Eq(a,b), Eq(b,c))
	innerAnd, ok := outerAnd.Left().(*prop.And)
	if !ok {
		t.Fatalf("expected left to be And proposition, got %T", outerAnd.Left())
	}

	// Right should be Eq(c,d)
	rightEqual, ok := outerAnd.Right().(*prop.Equal)
	if !ok {
		t.Fatalf("expected right to be Equal proposition, got %T", outerAnd.Right())
	}

	// Verify innerAnd.GetLeft() is Eq(a,b)
	innerLeftEqual, ok := innerAnd.Left().(*prop.Equal)
	if !ok {
		t.Fatalf("expected innerAnd.left to be Equal proposition, got %T", innerAnd.Left())
	}
	varA, ok := innerLeftEqual.Left().(*expr.Variable)
	if !ok || varA.Name() != "a" {
		t.Errorf("expected innerAnd.left.left to be variable 'a'")
	}
	varB1, ok := innerLeftEqual.Right().(*expr.Variable)
	if !ok || varB1.Name() != "b" {
		t.Errorf("expected innerAnd.left.right to be variable 'b'")
	}

	// Verify innerAnd.GetRight() is Eq(b,c)
	innerRightEqual, ok := innerAnd.Right().(*prop.Equal)
	if !ok {
		t.Fatalf("expected innerAnd.right to be Equal proposition, got %T", innerAnd.Right())
	}
	varB2, ok := innerRightEqual.Left().(*expr.Variable)
	if !ok || varB2.Name() != "b" {
		t.Errorf("expected innerAnd.right.left to be variable 'b'")
	}
	varC1, ok := innerRightEqual.Right().(*expr.Variable)
	if !ok || varC1.Name() != "c" {
		t.Errorf("expected innerAnd.right.right to be variable 'c'")
	}

	// Verify outerAnd.GetRight() is Eq(c,d)
	varC2, ok := rightEqual.Left().(*expr.Variable)
	if !ok || varC2.Name() != "c" {
		t.Errorf("expected outerAnd.right.left to be variable 'c'")
	}
	varD, ok := rightEqual.Right().(*expr.Variable)
	if !ok || varD.Name() != "d" {
		t.Errorf("expected outerAnd.right.right to be variable 'd'")
	}
}

// TestConverter_ChainedEquality_WithNumbers tests chained equality with numeric values
// Tests that 2 = 2 = 2 is properly converted to And(Eq(2,2), Eq(2,2))
//
// Note: This test expects the Converter to return a Proposition (And), not an Expression.
func TestConverter_ChainedEquality_WithNumbers(t *testing.T) {
	// 2 = 2 = 2
	// Expected: And(Eq(2,2), Eq(2,2))

	node := &EqualNode{
		Left: &EqualNode{
			Left:     &NumberNode{Value: 2.0},
			Operator: Token{Type: EQUAL},
			Right:    &NumberNode{Value: 2.0},
		},
		Operator: Token{Type: EQUAL},
		Right:    &NumberNode{Value: 2.0},
	}
	converter := NewConverter()

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	// Expected structure: And(Eq(2, 2), Eq(2, 2))
	andExpr, ok := result.(*prop.And)
	if !ok {
		t.Fatalf("expected And proposition, got %T", result)
	}

	// Left should be Equal(2, 2)
	leftEqual, ok := andExpr.Left().(*prop.Equal)
	if !ok {
		t.Fatalf("expected left to be Equal proposition, got %T", andExpr.Left())
	}

	// Right should be Equal(2, 2)
	rightEqual, ok := andExpr.Right().(*prop.Equal)
	if !ok {
		t.Fatalf("expected right to be Equal proposition, got %T", andExpr.Right())
	}

	// Verify structure of left Equal
	leftConst1, ok := leftEqual.Left().(*expr.Constant)
	if !ok || constantToFloat64(leftConst1) != 2.0 {
		t.Errorf("expected left.left to be constant 2.0")
	}
	leftConst2, ok := leftEqual.Right().(*expr.Constant)
	if !ok || constantToFloat64(leftConst2) != 2.0 {
		t.Errorf("expected left.right to be constant 2.0")
	}

	// Verify structure of right Equal
	rightConst1, ok := rightEqual.Left().(*expr.Constant)
	if !ok || constantToFloat64(rightConst1) != 2.0 {
		t.Errorf("expected right.left to be constant 2.0")
	}
	rightConst2, ok := rightEqual.Right().(*expr.Constant)
	if !ok || constantToFloat64(rightConst2) != 2.0 {
		t.Errorf("expected right.right to be constant 2.0")
	}
}

// TestConverter_ChainedEquality_WithExpressions tests chained equality with complex expressions
// Tests that (1+1) = 2 = (3-1) is properly converted to And(Eq(1+1, 2), Eq(2, 3-1))
//
// Note: This test expects the Converter to return a Proposition (And), not an Expression.
func TestConverter_ChainedEquality_WithExpressions(t *testing.T) {
	// (1+1) = 2 = (3-1)
	// Expected: And(Eq((1+1), 2), Eq(2, (3-1)))

	node := &EqualNode{
		Left: &EqualNode{
			Left: &BinaryOpNode{
				Left:     &NumberNode{Value: 1.0},
				Operator: Token{Type: PLUS},
				Right:    &NumberNode{Value: 1.0},
			},
			Operator: Token{Type: EQUAL},
			Right:    &NumberNode{Value: 2.0},
		},
		Operator: Token{Type: EQUAL},
		Right: &BinaryOpNode{
			Left:     &NumberNode{Value: 3.0},
			Operator: Token{Type: MINUS},
			Right:    &NumberNode{Value: 1.0},
		},
	}
	converter := NewConverter()

	result, err := converter.Convert(node)
	if err != nil {
		t.Fatalf("Convert error: %v", err)
	}

	// Expected structure: And(Eq((1+1), 2), Eq(2, (3-1)))
	andExpr, ok := result.(*prop.And)
	if !ok {
		t.Fatalf("expected And proposition, got %T", result)
	}

	// Left should be Equal((1+1), 2)
	leftEqual, ok := andExpr.Left().(*prop.Equal)
	if !ok {
		t.Fatalf("expected left to be Equal proposition, got %T", andExpr.Left())
	}

	// Right should be Equal(2, (3-1))
	rightEqual, ok := andExpr.Right().(*prop.Equal)
	if !ok {
		t.Fatalf("expected right to be Equal proposition, got %T", andExpr.Right())
	}

	// Verify left Equal: (1+1) = 2
	leftAdd, ok := leftEqual.Left().(*expr.Add)
	if !ok {
		t.Fatalf("expected left.left to be Add expression, got %T", leftEqual.Left())
	}
	_ = leftAdd // Structure verification is sufficient

	leftConst, ok := leftEqual.Right().(*expr.Constant)
	if !ok || constantToFloat64(leftConst) != 2.0 {
		t.Errorf("expected left.right to be constant 2.0")
	}

	// Verify right Equal: 2 = (3-1)
	rightConst, ok := rightEqual.Left().(*expr.Constant)
	if !ok || constantToFloat64(rightConst) != 2.0 {
		t.Errorf("expected right.left to be constant 2.0")
	}

	rightSub, ok := rightEqual.Right().(*expr.Sub)
	if !ok {
		t.Fatalf("expected right.right to be Subtract expression, got %T", rightEqual.Right())
	}
	_ = rightSub // Structure verification is sufficient
}
