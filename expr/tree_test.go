package expr

import (
	"math"
	"testing"
)

const floatEpsilon = 1e-9

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < floatEpsilon
}

func TestConstantEval(t *testing.T) {
	constant := &Constant{Value: NumberValue{Value: 5.0}}
	result, ok := constant.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 5.0 {
			t.Errorf("Expected 5.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestAddExpression(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 3.0}}
	right := &Constant{Value: NumberValue{Value: 4.0}}
	add := NewAddExpression(left, right)
	result, ok := add.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 7.0 {
			t.Errorf("Expected 7.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestSubtractExpression(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 10.0}}
	right := &Constant{Value: NumberValue{Value: 3.0}}
	sub := NewSubtractExpression(left, right)
	result, ok := sub.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 7.0 {
			t.Errorf("Expected 7.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestMultiplyExpression(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 6.0}}
	right := &Constant{Value: NumberValue{Value: 7.0}}
	mul := NewMultiplyExpression(left, right)
	result, ok := mul.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 42.0 {
			t.Errorf("Expected 42.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestDivideExpression(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 15.0}}
	right := &Constant{Value: NumberValue{Value: 3.0}}
	div := NewDivideExpression(left, right)
	result, ok := div.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 5.0 {
			t.Errorf("Expected 5.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestDivideByZero(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 10.0}}
	right := &Constant{Value: NumberValue{Value: 0.0}}
	div := NewDivideExpression(left, right)
	_, ok := div.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail due to division by zero")
	}
}

func TestAddExpressionChildren(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 3.0}}
	right := &Constant{Value: NumberValue{Value: 4.0}}
	add := NewAddExpression(left, right)
	children := add.Children()
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}
	if children[0] != left || children[1] != right {
		t.Errorf("Children do not match expected left and right")
	}
}

func TestConstantChildren(t *testing.T) {
	constant := &Constant{Value: NumberValue{Value: 5.0}}
	children := constant.Children()
	if len(children) != 0 {
		t.Errorf("Expected 0 children for Constant, got %d", len(children))
	}
}

func TestVariableEval(t *testing.T) {
	variable := &Variable{Name: "x"}
	_, ok := variable.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail for variable without context")
	}
}

func TestVariableChildren(t *testing.T) {
	variable := &Variable{Name: "x"}
	children := variable.Children()
	if len(children) != 0 {
		t.Errorf("Expected 0 children for Variable, got %d", len(children))
	}
}

func TestExpressionWithVariable(t *testing.T) {
	// x + 2 should fail to evaluate because x has no value
	left := &Variable{Name: "x"}
	right := &Constant{Value: NumberValue{Value: 2.0}}
	add := NewAddExpression(left, right)
	_, ok := add.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail due to variable")
	}
}

func TestPatternMatch(t *testing.T) {
	// Test case from requirement:
	// (x+3)(1+z) matches pattern x(y+z)
	// Expected bindings: x -> (x+3), y -> 1, z -> z

	// Build pattern: x(y+z)
	pattern := NewMultiplyExpression(
		NewVariable("x"),
		NewAddExpression(
			NewVariable("y"),
			NewVariable("z"),
		),
	)

	// Build expression: (x+3)(1+z)
	expr := NewMultiplyExpression(
		NewAddExpression(
			NewVariable("x"),
			NewConstant(3),
		),
		NewAddExpression(
			NewConstant(1),
			NewVariable("z"),
		),
	)

	bindings, ok := PatternMatch(pattern, expr)
	if !ok {
		t.Fatalf("Expected pattern match to succeed")
	}

	// Check binding for x: should be (x+3)
	xBinding, exists := bindings["x"]
	if !exists {
		t.Errorf("Expected binding for variable 'x'")
	}
	expectedX := NewAddExpression(
		NewVariable("x"),
		NewConstant(3),
	)
	if !expressionsEqual(xBinding, expectedX) {
		t.Errorf("Binding for 'x' does not match expected expression")
	}

	// Check binding for y: should be 1
	yBinding, exists := bindings["y"]
	if !exists {
		t.Errorf("Expected binding for variable 'y'")
	}
	expectedY := NewConstant(1)
	if !expressionsEqual(yBinding, expectedY) {
		t.Errorf("Binding for 'y' does not match expected expression")
	}

	// Check binding for z: should be z
	zBinding, exists := bindings["z"]
	if !exists {
		t.Errorf("Expected binding for variable 'z'")
	}
	expectedZ := NewVariable("z")
	if !expressionsEqual(zBinding, expectedZ) {
		t.Errorf("Binding for 'z' does not match expected expression")
	}
}

func TestPatternMatchWithRepeatedVariable(t *testing.T) {
	// Test that repeated variables must match the same expression
	// Pattern: x + x
	pattern := NewAddExpression(
		NewVariable("x"),
		NewVariable("x"),
	)

	// Expression: 2 + 2 (should match)
	expr1 := NewAddExpression(
		NewConstant(2),
		NewConstant(2),
	)
	bindings1, ok1 := PatternMatch(pattern, expr1)
	if !ok1 {
		t.Errorf("Expected pattern match to succeed for 2 + 2")
	}
	if xVal, exists := bindings1["x"]; !exists || !expressionsEqual(xVal, NewConstant(2)) {
		t.Errorf("Expected x to bind to 2")
	}

	// Expression: 2 + 3 (should NOT match)
	expr2 := NewAddExpression(
		NewConstant(2),
		NewConstant(3),
	)
	_, ok2 := PatternMatch(pattern, expr2)
	if ok2 {
		t.Errorf("Expected pattern match to fail for 2 + 3 (x must be consistent)")
	}
}

func TestPatternMatchWithConstant(t *testing.T) {
	// Pattern: x + 5
	pattern := NewAddExpression(
		NewVariable("x"),
		NewConstant(5),
	)

	// Expression: 3 + 5 (should match)
	expr1 := NewAddExpression(
		NewConstant(3),
		NewConstant(5),
	)
	bindings1, ok1 := PatternMatch(pattern, expr1)
	if !ok1 {
		t.Errorf("Expected pattern match to succeed")
	}
	if xVal, exists := bindings1["x"]; !exists || !expressionsEqual(xVal, NewConstant(3)) {
		t.Errorf("Expected x to bind to 3")
	}

	// Expression: 3 + 4 (should NOT match, constant differs)
	expr2 := NewAddExpression(
		NewConstant(3),
		NewConstant(4),
	)
	_, ok2 := PatternMatch(pattern, expr2)
	if ok2 {
		t.Errorf("Expected pattern match to fail (constant mismatch)")
	}
}

func TestPatternMatchTypeMismatch(t *testing.T) {
	// Pattern: x + y
	pattern := NewAddExpression(
		NewVariable("x"),
		NewVariable("y"),
	)

	// Expression: 2 * 3 (different operator, should NOT match)
	expr := NewMultiplyExpression(
		NewConstant(2),
		NewConstant(3),
	)

	_, ok := PatternMatch(pattern, expr)
	if ok {
		t.Errorf("Expected pattern match to fail (operator type mismatch)")
	}
}

func TestSubstitute(t *testing.T) {
	// Test case from requirement (inverse of PatternMatch):
	// Starting with: x(y+z)
	// Bindings: x -> (x+3), y -> 1, z -> z
	// Expected result: (x+3)(1+z)

	// Build expression: x(y+z)
	expr := NewMultiplyExpression(
		NewVariable("x"),
		NewAddExpression(
			NewVariable("y"),
			NewVariable("z"),
		),
	)

	// Build bindings
	bindings := map[string]Expression{
		"x": NewAddExpression(
			NewVariable("x"),
			NewConstant(3),
		),
		"y": NewConstant(1),
		"z": NewVariable("z"),
	}

	// Perform substitution
	result := Substitute(expr, bindings)

	// Expected result: (x+3)(1+z)
	expected := NewMultiplyExpression(
		NewAddExpression(
			NewVariable("x"),
			NewConstant(3),
		),
		NewAddExpression(
			NewConstant(1),
			NewVariable("z"),
		),
	)

	if !expressionsEqual(result, expected) {
		t.Errorf("Substitution result does not match expected expression")
	}
}

func TestSubstitutePartial(t *testing.T) {
	// Test partial substitution where not all variables are bound
	// Expression: x + y + z
	// Bindings: x -> 1, z -> 3 (y is not bound)
	// Expected: 1 + y + 3

	expr := NewAddExpression(
		NewAddExpression(
			NewVariable("x"),
			NewVariable("y"),
		),
		NewVariable("z"),
	)

	bindings := map[string]Expression{
		"x": NewConstant(1),
		"z": NewConstant(3),
	}

	result := Substitute(expr, bindings)

	expected := NewAddExpression(
		NewAddExpression(
			NewConstant(1),
			NewVariable("y"),
		),
		NewConstant(3),
	)

	if !expressionsEqual(result, expected) {
		t.Errorf("Partial substitution result does not match expected expression")
	}
}

func TestSubstituteNoBindings(t *testing.T) {
	// Test substitution with empty bindings
	// Expression: x + y
	// Bindings: {} (empty)
	// Expected: x + y (unchanged)

	expr := NewAddExpression(
		NewVariable("x"),
		NewVariable("y"),
	)

	bindings := map[string]Expression{}

	result := Substitute(expr, bindings)

	if !expressionsEqual(result, expr) {
		t.Errorf("Expression should remain unchanged with empty bindings")
	}
}

func TestSubstituteWithConstants(t *testing.T) {
	// Test that constants are preserved during substitution
	// Expression: 2 + x * 3
	// Bindings: x -> 5
	// Expected: 2 + 5 * 3

	expr := NewAddExpression(
		NewConstant(2),
		NewMultiplyExpression(
			NewVariable("x"),
			NewConstant(3),
		),
	)

	bindings := map[string]Expression{
		"x": NewConstant(5),
	}

	result := Substitute(expr, bindings)

	expected := NewAddExpression(
		NewConstant(2),
		NewMultiplyExpression(
			NewConstant(5),
			NewConstant(3),
		),
	)

	if !expressionsEqual(result, expected) {
		t.Errorf("Substitution with constants does not match expected expression")
	}
}

func TestPatternMatchAndSubstituteRoundTrip(t *testing.T) {
	// Test that PatternMatch and Substitute are inverse operations
	// 1. Match (x+3)(1+z) against pattern x(y+z) to get bindings
	// 2. Substitute the bindings back into the pattern
	// 3. Should get back the original expression

	pattern := NewMultiplyExpression(
		NewVariable("x"),
		NewAddExpression(
			NewVariable("y"),
			NewVariable("z"),
		),
	)

	original := NewMultiplyExpression(
		NewAddExpression(
			NewVariable("x"),
			NewConstant(3),
		),
		NewAddExpression(
			NewConstant(1),
			NewVariable("z"),
		),
	)

	// Step 1: Pattern match
	bindings, ok := PatternMatch(pattern, original)
	if !ok {
		t.Fatalf("Pattern match failed")
	}

	// Step 2: Substitute back
	result := Substitute(pattern, bindings)

	// Step 3: Compare
	if !expressionsEqual(result, original) {
		t.Errorf("Round-trip failed: result does not match original expression")
	}
}

// ===== Power Expression Tests =====

func TestPowerExpressionBasic(t *testing.T) {
	// 2^3 = 8
	base := &Constant{Value: NumberValue{Value: 2.0}}
	exponent := &Constant{Value: NumberValue{Value: 3.0}}
	pow := NewPowerExpression(base, exponent)
	result, ok := pow.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 8.0 {
			t.Errorf("Expected 8.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestPowerExpressionZeroExponent(t *testing.T) {
	// 5^0 = 1
	base := &Constant{Value: NumberValue{Value: 5.0}}
	exponent := &Constant{Value: NumberValue{Value: 0.0}}
	pow := NewPowerExpression(base, exponent)
	result, ok := pow.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 1.0 {
			t.Errorf("Expected 1.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestPowerExpressionOneExponent(t *testing.T) {
	// 5^1 = 5
	base := &Constant{Value: NumberValue{Value: 5.0}}
	exponent := &Constant{Value: NumberValue{Value: 1.0}}
	pow := NewPowerExpression(base, exponent)
	result, ok := pow.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 5.0 {
			t.Errorf("Expected 5.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestPowerExpressionNegativeExponent(t *testing.T) {
	// 2^(-2) = 0.25
	base := &Constant{Value: NumberValue{Value: 2.0}}
	exponent := &Constant{Value: NumberValue{Value: -2.0}}
	pow := NewPowerExpression(base, exponent)
	result, ok := pow.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 0.25 {
			t.Errorf("Expected 0.25, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestPowerExpressionFractionalExponent(t *testing.T) {
	// 4^0.5 = 2
	base := &Constant{Value: NumberValue{Value: 4.0}}
	exponent := &Constant{Value: NumberValue{Value: 0.5}}
	pow := NewPowerExpression(base, exponent)
	result, ok := pow.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 2.0 {
			t.Errorf("Expected 2.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestPowerExpressionNested(t *testing.T) {
	// (2^2)^3 = 4^3 = 64
	inner := NewPowerExpression(
		&Constant{Value: NumberValue{Value: 2.0}},
		&Constant{Value: NumberValue{Value: 2.0}},
	)
	outer := NewPowerExpression(
		inner,
		&Constant{Value: NumberValue{Value: 3.0}},
	)
	result, ok := outer.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 64.0 {
			t.Errorf("Expected 64.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestPowerExpressionZeroToZero(t *testing.T) {
	// 0^0 is mathematically undefined, but many implementations define it as 1
	base := &Constant{Value: NumberValue{Value: 0.0}}
	exponent := &Constant{Value: NumberValue{Value: 0.0}}
	pow := NewPowerExpression(base, exponent)
	result, ok := pow.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 1.0 {
			t.Errorf("Expected 1.0 for 0^0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestPowerExpressionChildren(t *testing.T) {
	base := &Constant{Value: NumberValue{Value: 2.0}}
	exponent := &Constant{Value: NumberValue{Value: 3.0}}
	pow := NewPowerExpression(base, exponent)
	children := pow.Children()
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}
	if children[0] != base || children[1] != exponent {
		t.Errorf("Children do not match expected base and exponent")
	}
}

// ===== Square Root Expression Tests =====

func TestSqrtExpressionBasic(t *testing.T) {
	// sqrt(4) = 2
	arg := &Constant{Value: NumberValue{Value: 4.0}}
	sqrt := NewSqrtExpression(arg)
	result, ok := sqrt.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 2.0 {
			t.Errorf("Expected 2.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestSqrtExpressionZero(t *testing.T) {
	// sqrt(0) = 0
	arg := &Constant{Value: NumberValue{Value: 0.0}}
	sqrt := NewSqrtExpression(arg)
	result, ok := sqrt.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 0.0 {
			t.Errorf("Expected 0.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestSqrtExpressionOne(t *testing.T) {
	// sqrt(1) = 1
	arg := &Constant{Value: NumberValue{Value: 1.0}}
	sqrt := NewSqrtExpression(arg)
	result, ok := sqrt.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 1.0 {
			t.Errorf("Expected 1.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestSqrtExpressionDecimal(t *testing.T) {
	// sqrt(2.25) = 1.5
	arg := &Constant{Value: NumberValue{Value: 2.25}}
	sqrt := NewSqrtExpression(arg)
	result, ok := sqrt.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 1.5 {
			t.Errorf("Expected 1.5, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestSqrtExpressionNegative(t *testing.T) {
	// sqrt(-1) should fail (no complex number support)
	arg := &Constant{Value: NumberValue{Value: -1.0}}
	sqrt := NewSqrtExpression(arg)
	_, ok := sqrt.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail for negative number")
	}
}

func TestSqrtExpressionNested(t *testing.T) {
	// sqrt(sqrt(16)) = sqrt(4) = 2
	inner := NewSqrtExpression(
		&Constant{Value: NumberValue{Value: 16.0}},
	)
	outer := NewSqrtExpression(inner)
	result, ok := outer.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 2.0 {
			t.Errorf("Expected 2.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestSqrtExpressionChildren(t *testing.T) {
	arg := &Constant{Value: NumberValue{Value: 4.0}}
	sqrt := NewSqrtExpression(arg)
	children := sqrt.Children()
	if len(children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(children))
	}
	if children[0] != arg {
		t.Errorf("Child does not match expected argument")
	}
}

func TestSqrtExpressionWithComplexExpression(t *testing.T) {
	// sqrt(3^2 + 4^2) = sqrt(9 + 16) = sqrt(25) = 5
	// This tests the Pythagorean theorem
	pow1 := NewPowerExpression(
		&Constant{Value: NumberValue{Value: 3.0}},
		&Constant{Value: NumberValue{Value: 2.0}},
	)
	pow2 := NewPowerExpression(
		&Constant{Value: NumberValue{Value: 4.0}},
		&Constant{Value: NumberValue{Value: 2.0}},
	)
	sum := NewAddExpression(pow1, pow2)
	sqrt := NewSqrtExpression(sum)
	result, ok := sqrt.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 5.0 {
			t.Errorf("Expected 5.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

// ===== Nth Root Expression Tests =====

func TestCubeRoot(t *testing.T) {
	// ∛8 = 2
	arg := &Constant{Value: NumberValue{Value: 8.0}}
	cubeRoot := NewNthRootExpression(arg, 3)
	result, ok := cubeRoot.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if !almostEqual(num.Value, 2.0) {
			t.Errorf("Expected 2.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestFourthRoot(t *testing.T) {
	// ⁴√16 = 2
	arg := &Constant{Value: NumberValue{Value: 16.0}}
	fourthRoot := NewNthRootExpression(arg, 4)
	result, ok := fourthRoot.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 2.0 {
			t.Errorf("Expected 2.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestCubeRootOfNegative(t *testing.T) {
	// ∛(-8) = -2 (odd roots of negative numbers are allowed)
	arg := &Constant{Value: NumberValue{Value: -8.0}}
	cubeRoot := NewNthRootExpression(arg, 3)
	result, ok := cubeRoot.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if !almostEqual(num.Value, -2.0) {
			t.Errorf("Expected -2.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestEvenRootOfNegative(t *testing.T) {
	// ⁴√(-16) should fail (even roots of negative numbers are not allowed)
	arg := &Constant{Value: NumberValue{Value: -16.0}}
	fourthRoot := NewNthRootExpression(arg, 4)
	_, ok := fourthRoot.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail for even root of negative number")
	}
}

func TestFifthRoot(t *testing.T) {
	// ⁵√32 = 2
	arg := &Constant{Value: NumberValue{Value: 32.0}}
	fifthRoot := NewNthRootExpression(arg, 5)
	result, ok := fifthRoot.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if num, ok := result.(*NumberValue); ok {
		if num.Value != 2.0 {
			t.Errorf("Expected 2.0, got %f", num.Value)
		}
	} else {
		t.Errorf("Expected NumberValue")
	}
}

func TestNthRootChildren(t *testing.T) {
	arg := &Constant{Value: NumberValue{Value: 8.0}}
	cubeRoot := NewNthRootExpression(arg, 3)
	children := cubeRoot.Children()
	if len(children) != 1 {
		t.Errorf("Expected 1 child, got %d", len(children))
	}
	if children[0] != arg {
		t.Errorf("Child does not match expected argument")
	}
}

func TestNthRootNValue(t *testing.T) {
	// Test that N is correctly set
	cubeRoot := NewNthRootExpression(&Constant{Value: NumberValue{Value: 8.0}}, 3)
	if cubeRoot.N != 3.0 {
		t.Errorf("Expected N to be 3.0, got %f", cubeRoot.N)
	}

	// Test that NewSqrtExpression sets N to 2
	sqrtExpr := NewSqrtExpression(&Constant{Value: NumberValue{Value: 4.0}})
	if sqrtExpr.N != 2.0 {
		t.Errorf("Expected N to be 2.0 for NewSqrtExpression, got %f", sqrtExpr.N)
	}
}

// ===== BoolValue Tests =====

func TestBoolValueEval(t *testing.T) {
	// Test BoolValue evaluation
	boolTrue := &BoolValue{Value: true}
	resultTrue, okTrue := boolTrue.Eval()
	if !okTrue {
		t.Errorf("Expected evaluation to succeed for BoolValue")
	}
	if b, ok := resultTrue.(*BoolValue); ok {
		if !b.Value {
			t.Errorf("Expected true, got false")
		}
	} else {
		t.Errorf("Expected BoolValue")
	}

	boolFalse := &BoolValue{Value: false}
	resultFalse, okFalse := boolFalse.Eval()
	if !okFalse {
		t.Errorf("Expected evaluation to succeed for BoolValue")
	}
	if b, ok := resultFalse.(*BoolValue); ok {
		if b.Value {
			t.Errorf("Expected false, got true")
		}
	} else {
		t.Errorf("Expected BoolValue")
	}
}

// ===== EqualExpression Tests =====

func TestEqualExpressionTrue(t *testing.T) {
	// 5 = 5 should be true
	left := &Constant{Value: NumberValue{Value: 5.0}}
	right := &Constant{Value: NumberValue{Value: 5.0}}
	equal := NewEqualExpression(left, right)
	result, ok := equal.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if b, ok := result.(*BoolValue); ok {
		if !b.Value {
			t.Errorf("Expected true for 5 = 5")
		}
	} else {
		t.Errorf("Expected BoolValue, got %T", result)
	}
}

func TestEqualExpressionFalse(t *testing.T) {
	// 2 = 3 should be false
	left := &Constant{Value: NumberValue{Value: 2.0}}
	right := &Constant{Value: NumberValue{Value: 3.0}}
	equal := NewEqualExpression(left, right)
	result, ok := equal.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if b, ok := result.(*BoolValue); ok {
		if b.Value {
			t.Errorf("Expected false for 2 = 3")
		}
	} else {
		t.Errorf("Expected BoolValue, got %T", result)
	}
}

func TestEqualExpressionMathematical(t *testing.T) {
	// 2 + 3 = 5 (数学的に正しい等式)
	left := NewAddExpression(
		NewConstant(2),
		NewConstant(3),
	)
	right := NewConstant(5)
	equal := NewEqualExpression(left, right)

	result, ok := equal.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}

	boolResult, ok := result.(*BoolValue)
	if !ok {
		t.Fatalf("Expected BoolValue, got %T", result)
	}

	if !boolResult.Value {
		t.Errorf("Expected true for 2+3=5")
	}
}

func TestEqualExpressionFloatingPointTolerance(t *testing.T) {
	// 0.1 + 0.2 = 0.3 (浮動小数点誤差を考慮)
	left := NewAddExpression(
		&Constant{Value: NumberValue{Value: 0.1}},
		&Constant{Value: NumberValue{Value: 0.2}},
	)
	right := &Constant{Value: NumberValue{Value: 0.3}}
	equal := NewEqualExpression(left, right)

	result, ok := equal.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}

	boolResult, ok := result.(*BoolValue)
	if !ok {
		t.Fatalf("Expected BoolValue, got %T", result)
	}

	if !boolResult.Value {
		t.Errorf("Expected true for 0.1+0.2=0.3 (with floating point tolerance)")
	}
}

func TestEqualExpressionComplex(t *testing.T) {
	// 2 + 3 = 1 + 4 (both sides evaluate to 5)
	left := NewAddExpression(
		NewConstant(2),
		NewConstant(3),
	)
	right := NewAddExpression(
		NewConstant(1),
		NewConstant(4),
	)
	equal := NewEqualExpression(left, right)

	result, ok := equal.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}

	boolResult, ok := result.(*BoolValue)
	if !ok {
		t.Fatalf("Expected BoolValue, got %T", result)
	}

	if !boolResult.Value {
		t.Errorf("Expected true for (2+3)=(1+4)")
	}
}

func TestEqualExpressionZero(t *testing.T) {
	// 0 = 0 should be true
	left := &Constant{Value: NumberValue{Value: 0.0}}
	right := &Constant{Value: NumberValue{Value: 0.0}}
	equal := NewEqualExpression(left, right)
	result, ok := equal.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if b, ok := result.(*BoolValue); ok {
		if !b.Value {
			t.Errorf("Expected true for 0 = 0")
		}
	} else {
		t.Errorf("Expected BoolValue")
	}
}

func TestEqualExpressionNegative(t *testing.T) {
	// -5 = -5 should be true
	left := &Constant{Value: NumberValue{Value: -5.0}}
	right := &Constant{Value: NumberValue{Value: -5.0}}
	equal := NewEqualExpression(left, right)
	result, ok := equal.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}
	if b, ok := result.(*BoolValue); ok {
		if !b.Value {
			t.Errorf("Expected true for -5 = -5")
		}
	} else {
		t.Errorf("Expected BoolValue")
	}
}

func TestEqualExpressionEvaluationFailure(t *testing.T) {
	// x = 5 should fail because x cannot be evaluated
	left := &Variable{Name: "x"}
	right := &Constant{Value: NumberValue{Value: 5.0}}
	equal := NewEqualExpression(left, right)
	_, ok := equal.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail due to variable")
	}
}

func TestEqualExpressionBothSidesFail(t *testing.T) {
	// x = y should fail because both variables cannot be evaluated
	left := &Variable{Name: "x"}
	right := &Variable{Name: "y"}
	equal := NewEqualExpression(left, right)
	_, ok := equal.Eval()
	if ok {
		t.Errorf("Expected evaluation to fail due to variables")
	}
}

func TestEqualExpressionChildren(t *testing.T) {
	left := &Constant{Value: NumberValue{Value: 2.0}}
	right := &Constant{Value: NumberValue{Value: 3.0}}
	equal := NewEqualExpression(left, right)
	children := equal.Children()
	if len(children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(children))
	}
	if children[0] != left || children[1] != right {
		t.Errorf("Children do not match expected left and right")
	}
}

func TestEqualExpressionBoolValues(t *testing.T) {
	// (2 = 2) = (3 = 3) should be true (both sides are true)
	leftEqual := NewEqualExpression(
		NewConstant(2),
		NewConstant(2),
	)
	rightEqual := NewEqualExpression(
		NewConstant(3),
		NewConstant(3),
	)
	outerEqual := NewEqualExpression(leftEqual, rightEqual)

	result, ok := outerEqual.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}

	boolResult, ok := result.(*BoolValue)
	if !ok {
		t.Fatalf("Expected BoolValue, got %T", result)
	}

	if !boolResult.Value {
		t.Errorf("Expected true for (2=2)=(3=3)")
	}
}

func TestEqualExpressionBoolValuesFalse(t *testing.T) {
	// (2 = 2) = (3 = 4) should be false (true != false)
	leftEqual := NewEqualExpression(
		NewConstant(2),
		NewConstant(2),
	)
	rightEqual := NewEqualExpression(
		NewConstant(3),
		NewConstant(4),
	)
	outerEqual := NewEqualExpression(leftEqual, rightEqual)

	result, ok := outerEqual.Eval()
	if !ok {
		t.Errorf("Expected evaluation to succeed")
	}

	boolResult, ok := result.(*BoolValue)
	if !ok {
		t.Fatalf("Expected BoolValue, got %T", result)
	}

	if boolResult.Value {
		t.Errorf("Expected false for (2=2)=(3=4)")
	}
}

func TestPatternMatchEqual(t *testing.T) {
	// Pattern: x = y
	pattern := NewEqualExpression(
		NewVariable("x"),
		NewVariable("y"),
	)

	// Expression: 2 + 3 = 5
	expr := NewEqualExpression(
		NewAddExpression(
			NewConstant(2),
			NewConstant(3),
		),
		NewConstant(5),
	)

	bindings, ok := PatternMatch(pattern, expr)
	if !ok {
		t.Fatalf("Expected pattern match to succeed")
	}

	// Check binding for x: should be (2 + 3)
	xBinding, exists := bindings["x"]
	if !exists {
		t.Errorf("Expected binding for variable 'x'")
	}
	expectedX := NewAddExpression(
		NewConstant(2),
		NewConstant(3),
	)
	if !expressionsEqual(xBinding, expectedX) {
		t.Errorf("Binding for 'x' does not match expected expression")
	}

	// Check binding for y: should be 5
	yBinding, exists := bindings["y"]
	if !exists {
		t.Errorf("Expected binding for variable 'y'")
	}
	expectedY := NewConstant(5)
	if !expressionsEqual(yBinding, expectedY) {
		t.Errorf("Binding for 'y' does not match expected expression")
	}
}

func TestSubstituteEqual(t *testing.T) {
	// Expression: x = y
	expr := NewEqualExpression(
		NewVariable("x"),
		NewVariable("y"),
	)

	// Bindings: x -> 2 + 3, y -> 5
	bindings := map[string]Expression{
		"x": NewAddExpression(
			NewConstant(2),
			NewConstant(3),
		),
		"y": NewConstant(5),
	}

	// Perform substitution
	result := Substitute(expr, bindings)

	// Expected result: (2 + 3) = 5
	expected := NewEqualExpression(
		NewAddExpression(
			NewConstant(2),
			NewConstant(3),
		),
		NewConstant(5),
	)

	if !expressionsEqual(result, expected) {
		t.Errorf("Substitution result does not match expected expression")
	}
}
