package expr

import (
	"testing"
)

// TestSplitToFactors tests the SplitToFactors function that decomposes an expression into multiplicative factors
func TestSplitToFactors(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected []Expression
	}{
		{
			name:     "single constant",
			expr:     NewConstant(5),
			expected: []Expression{NewConstant(5)},
		},
		{
			name:     "single variable",
			expr:     NewVariable("x"),
			expected: []Expression{NewVariable("x")},
		},
		{
			name: "simple multiplication: 2 * 3",
			expr: NewMultiply(
				NewConstant(2),
				NewConstant(3),
			),
			expected: []Expression{
				NewConstant(2),
				NewConstant(3),
			},
		},
		{
			name: "nested multiplication: (2 * 3) * 4",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			),
			expected: []Expression{
				NewConstant(2),
				NewConstant(3),
				NewConstant(4),
			},
		},
		{
			name: "deeply nested: ((2 * 3) * 4) * 5",
			expr: NewMultiply(
				NewMultiply(
					NewMultiply(
						NewConstant(2),
						NewConstant(3),
					),
					NewConstant(4),
				),
				NewConstant(5),
			),
			expected: []Expression{
				NewConstant(2),
				NewConstant(3),
				NewConstant(4),
				NewConstant(5),
			},
		},
		{
			name: "variables and constants: 2 * x * y",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewVariable("y"),
			),
			expected: []Expression{
				NewConstant(2),
				NewVariable("x"),
				NewVariable("y"),
			},
		},
		{
			name: "complex monomial: 3 * x * y * z",
			expr: NewMultiply(
				NewMultiply(
					NewMultiply(
						NewConstant(3),
						NewVariable("x"),
					),
					NewVariable("y"),
				),
				NewVariable("z"),
			),
			expected: []Expression{
				NewConstant(3),
				NewVariable("x"),
				NewVariable("y"),
				NewVariable("z"),
			},
		},
		{
			name: "addition only: 2 + 3",
			expr: NewAdd(
				NewConstant(2),
				NewConstant(3),
			),
			expected: []Expression{
				NewAdd(
					NewConstant(2),
					NewConstant(3),
				),
			},
		},
		{
			name: "division only: 10 / 2",
			expr: NewDivide(
				NewConstant(10),
				NewConstant(2),
			),
			expected: []Expression{
				NewDivide(
					NewConstant(10),
					NewConstant(2),
				),
			},
		},
		{
			name: "subtraction only: 5 - 3",
			expr: NewSubtract(
				NewConstant(5),
				NewConstant(3),
			),
			expected: []Expression{
				NewSubtract(
					NewConstant(5),
					NewConstant(3),
				),
			},
		},
		{
			name: "multiplication with addition: (2 + 3) * 4",
			expr: NewMultiply(
				NewAdd(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			),
			expected: []Expression{
				NewAdd(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			},
		},
		{
			name: "power expression: x^2",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(2),
			),
			expected: []Expression{
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
			},
		},
		{
			name: "multiplication with power: 3 * x^2",
			expr: NewMultiply(
				NewConstant(3),
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
			),
			expected: []Expression{
				NewConstant(3),
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
			},
		},
		{
			name: "sqrt expression: sqrt(x)",
			expr: NewSqrt(
				NewVariable("x"),
			),
			expected: []Expression{
				NewSqrt(
					NewVariable("x"),
				),
			},
		},
		{
			name: "multiplication with sqrt: 2 * sqrt(x) * y",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewSqrt(
						NewVariable("x"),
					),
				),
				NewVariable("y"),
			),
			expected: []Expression{
				NewConstant(2),
				NewSqrt(
					NewVariable("x"),
				),
				NewVariable("y"),
			},
		},
		{
			name: "nth root expression: cbrt(8)",
			expr: NewNthRoot(
				NewConstant(8),
				3,
			),
			expected: []Expression{
				NewNthRoot(
					NewConstant(8),
					3,
				),
			},
		},
		{
			name: "nested power: (x^2)^3",
			expr: NewPower(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected: []Expression{
				NewPower(
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
					NewConstant(3),
				),
			},
		},
		{
			name: "power of sqrt: (sqrt(x))^2",
			expr: NewPower(
				NewSqrt(
					NewVariable("x"),
				),
				NewConstant(2),
			),
			expected: []Expression{
				NewPower(
					NewSqrt(
						NewVariable("x"),
					),
					NewConstant(2),
				),
			},
		},
		{
			name: "sqrt of power: sqrt(x^2)",
			expr: NewSqrt(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
			),
			expected: []Expression{
				NewSqrt(
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
				),
			},
		},
		{
			name: "multiple powers: x^2 * y^3 * z^4",
			expr: NewMultiply(
				NewMultiply(
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
					NewPower(
						NewVariable("y"),
						NewConstant(3),
					),
				),
				NewPower(
					NewVariable("z"),
					NewConstant(4),
				),
			),
			expected: []Expression{
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewPower(
					NewVariable("y"),
					NewConstant(3),
				),
				NewPower(
					NewVariable("z"),
					NewConstant(4),
				),
			},
		},
		{
			name: "different roots: sqrt(x) * cbrt(y)",
			expr: NewMultiply(
				NewSqrt(
					NewVariable("x"),
				),
				NewNthRoot(
					NewVariable("y"),
					3,
				),
			),
			expected: []Expression{
				NewSqrt(
					NewVariable("x"),
				),
				NewNthRoot(
					NewVariable("y"),
					3,
				),
			},
		},
		{
			name: "complex: 2 * x^2 * sqrt(y) * z",
			expr: NewMultiply(
				NewMultiply(
					NewMultiply(
						NewConstant(2),
						NewPower(
							NewVariable("x"),
							NewConstant(2),
						),
					),
					NewSqrt(
						NewVariable("y"),
					),
				),
				NewVariable("z"),
			),
			expected: []Expression{
				NewConstant(2),
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewSqrt(
					NewVariable("y"),
				),
				NewVariable("z"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitToFactors(tt.expr)

			if len(result) != len(tt.expected) {
				t.Errorf("SplitToFactors() returned %d factors, expected %d", len(result), len(tt.expected))
				return
			}

			for i, factor := range result {
				if !expressionsEqual(factor, tt.expected[i]) {
					t.Errorf("SplitToFactors() factor[%d] mismatch", i)
				}
			}
		})
	}
}

// TestCombineFactors tests the CombineFactors function that combines a slice of expressions using multiplication
func TestCombineFactors(t *testing.T) {
	tests := []struct {
		name     string
		factors  []Expression
		expected Expression
	}{
		{
			name:     "single factor: 5",
			factors:  []Expression{NewConstant(5)},
			expected: NewConstant(5),
		},
		{
			name: "two factors: 2 * 3",
			factors: []Expression{
				NewConstant(2),
				NewConstant(3),
			},
			expected: NewMultiply(
				NewConstant(2),
				NewConstant(3),
			),
		},
		{
			name: "three factors: 2 * 3 * 4",
			factors: []Expression{
				NewConstant(2),
				NewConstant(3),
				NewConstant(4),
			},
			expected: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			),
		},
		{
			name: "four factors: 2 * 3 * 4 * 5",
			factors: []Expression{
				NewConstant(2),
				NewConstant(3),
				NewConstant(4),
				NewConstant(5),
			},
			expected: NewMultiply(
				NewMultiply(
					NewMultiply(
						NewConstant(2),
						NewConstant(3),
					),
					NewConstant(4),
				),
				NewConstant(5),
			),
		},
		{
			name: "variables and constants: 2 * x * y",
			factors: []Expression{
				NewConstant(2),
				NewVariable("x"),
				NewVariable("y"),
			},
			expected: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewVariable("y"),
			),
		},
		{
			name: "complex factors with addition: (2 + 3) * x",
			factors: []Expression{
				NewAdd(
					NewConstant(2),
					NewConstant(3),
				),
				NewVariable("x"),
			},
			expected: NewMultiply(
				NewAdd(
					NewConstant(2),
					NewConstant(3),
				),
				NewVariable("x"),
			),
		},
		{
			name:     "empty slice",
			factors:  []Expression{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CombineFactors(tt.factors)

			// Handle nil case
			if tt.expected == nil {
				if result != nil {
					t.Errorf("CombineFactors() expected nil, got non-nil")
				}
				return
			}

			if result == nil {
				t.Errorf("CombineFactors() returned nil, expected non-nil")
				return
			}

			if !expressionsEqual(result, tt.expected) {
				t.Errorf("CombineFactors() result mismatch")
			}
		})
	}
}

// TestRoundTripSplitCombineFactors tests that splitting and combining factors is an identity operation
func TestRoundTripSplitCombineFactors(t *testing.T) {
	tests := []struct {
		name string
		expr Expression
	}{
		{
			name: "simple multiplication: 2 * 3",
			expr: NewMultiply(
				NewConstant(2),
				NewConstant(3),
			),
		},
		{
			name: "nested multiplication: (2 * 3) * 4",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			),
		},
		{
			name: "variables and constants: 2 * x * y",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewVariable("y"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Split and then combine
			factors := SplitToFactors(tt.expr)
			result := CombineFactors(factors)

			// Check if the result is structurally equal to the original
			if !expressionsEqual(result, tt.expr) {
				t.Errorf("Round trip failed: split and combine did not produce equivalent expression")
			}
		})
	}
}

// TestCountFactors tests the CountFactors function that counts the number of multiplicative factors
func TestCountFactors(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected int
	}{
		{
			name:     "single constant",
			expr:     NewConstant(5),
			expected: 1,
		},
		{
			name:     "single variable",
			expr:     NewVariable("x"),
			expected: 1,
		},
		{
			name: "simple multiplication: 2 * 3",
			expr: NewMultiply(
				NewConstant(2),
				NewConstant(3),
			),
			expected: 2,
		},
		{
			name: "nested multiplication: (2 * 3) * 4",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			),
			expected: 3,
		},
		{
			name: "deeply nested: ((2 * 3) * 4) * 5",
			expr: NewMultiply(
				NewMultiply(
					NewMultiply(
						NewConstant(2),
						NewConstant(3),
					),
					NewConstant(4),
				),
				NewConstant(5),
			),
			expected: 4,
		},
		{
			name: "addition only: 2 + 3",
			expr: NewAdd(
				NewConstant(2),
				NewConstant(3),
			),
			expected: 1,
		},
		{
			name: "variables: x * y * z",
			expr: NewMultiply(
				NewMultiply(
					NewVariable("x"),
					NewVariable("y"),
				),
				NewVariable("z"),
			),
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountFactors(tt.expr)

			if result != tt.expected {
				t.Errorf("CountFactors() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

// TestIsMonomial tests whether an expression is a valid monomial (no addition or subtraction)
func TestIsMonomial(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected bool
	}{
		{
			name:     "constant",
			expr:     NewConstant(5),
			expected: true,
		},
		{
			name:     "variable",
			expr:     NewVariable("x"),
			expected: true,
		},
		{
			name: "simple multiplication: 2 * x",
			expr: NewMultiply(
				NewConstant(2),
				NewVariable("x"),
			),
			expected: true,
		},
		{
			name: "nested multiplication: 2 * x * y",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewVariable("y"),
			),
			expected: true,
		},
		{
			name: "addition (not a monomial)",
			expr: NewAdd(
				NewConstant(2),
				NewConstant(3),
			),
			expected: false,
		},
		{
			name: "subtraction (not a monomial)",
			expr: NewSubtract(
				NewConstant(5),
				NewConstant(3),
			),
			expected: false,
		},
		{
			name: "multiplication containing addition (not a monomial)",
			expr: NewMultiply(
				NewAdd(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected: false,
		},
		{
			name: "division with constant denominator",
			expr: NewDivide(
				NewVariable("x"),
				NewConstant(2),
			),
			expected: true,
		},
		{
			name: "division with variable denominator (not a monomial in strict sense)",
			expr: NewDivide(
				NewConstant(2),
				NewVariable("x"),
			),
			expected: false,
		},
		{
			name: "complex monomial: 3 * x * y * z",
			expr: NewMultiply(
				NewMultiply(
					NewMultiply(
						NewConstant(3),
						NewVariable("x"),
					),
					NewVariable("y"),
				),
				NewVariable("z"),
			),
			expected: true,
		},
		{
			name: "power expression: x^2",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(2),
			),
			expected: true,
		},
		{
			name: "power with coefficient: 3 * x^2",
			expr: NewMultiply(
				NewConstant(3),
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
			),
			expected: true,
		},
		{
			name: "square root: sqrt(x)",
			expr: NewSqrt(
				NewVariable("x"),
			),
			expected: false,
		},
		{
			name: "nth root: cbrt(x)",
			expr: NewNthRoot(
				NewVariable("x"),
				3,
			),
			expected: false,
		},
		{
			name: "coefficient with sqrt: 2 * sqrt(x)",
			expr: NewMultiply(
				NewConstant(2),
				NewSqrt(
					NewVariable("x"),
				),
			),
			expected: false,
		},
		{
			name: "power and variable: x^2 * y",
			expr: NewMultiply(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewVariable("y"),
			),
			expected: true,
		},
		{
			name: "nested power: (x^2)^3",
			expr: NewPower(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected: true,
		},
		{
			name: "power of sqrt: (sqrt(x))^2",
			expr: NewPower(
				NewSqrt(
					NewVariable("x"),
				),
				NewConstant(2),
			),
			expected: false,
		},
		{
			name: "sqrt of power: sqrt(x^2)",
			expr: NewSqrt(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
			),
			expected: false,
		},
		{
			name: "multiple roots: sqrt(x) * sqrt(y)",
			expr: NewMultiply(
				NewSqrt(
					NewVariable("x"),
				),
				NewSqrt(
					NewVariable("y"),
				),
			),
			expected: false,
		},
		{
			name: "different root degrees: sqrt(x) * cbrt(x)",
			expr: NewMultiply(
				NewSqrt(
					NewVariable("x"),
				),
				NewNthRoot(
					NewVariable("x"),
					3,
				),
			),
			expected: false,
		},
		{
			name: "complex monomial: 2 * x^2 * sqrt(y) * z",
			expr: NewMultiply(
				NewMultiply(
					NewMultiply(
						NewConstant(2),
						NewPower(
							NewVariable("x"),
							NewConstant(2),
						),
					),
					NewSqrt(
						NewVariable("y"),
					),
				),
				NewVariable("z"),
			),
			expected: false,
		},
		{
			name: "power with zero exponent: x^0",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(0),
			),
			expected: false,
		},
		{
			name: "power with one exponent: x^1",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(1),
			),
			expected: true,
		},
		{
			name: "negative power: x^(-2) (not a monomial in strict sense)",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(-2),
			),
			expected: false,
		},
		{
			name: "fractional power: x^(1/2)",
			expr: NewPower(
				NewVariable("x"),
				NewDivide(
					NewConstant(1),
					NewConstant(2),
				),
			),
			expected: false,
		},
		{
			name: "power with variable exponent: x^y (not a standard monomial)",
			expr: NewPower(
				NewVariable("x"),
				NewVariable("y"),
			),
			expected: false,
		},
		{
			name: "sqrt of complex expression: sqrt(x + 1) (not a monomial)",
			expr: NewSqrt(
				NewAdd(
					NewVariable("x"),
					NewConstant(1),
				),
			),
			expected: false,
		},
		{
			name: "power of sum: (x + y)^2 (not a monomial)",
			expr: NewPower(
				NewAdd(
					NewVariable("x"),
					NewVariable("y"),
				),
				NewConstant(2),
			),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsMonomial(tt.expr)

			if result != tt.expected {
				t.Errorf("IsMonomial() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestMapFactors tests the MapFactors function that applies a transformation to each factor
func TestMapFactors(t *testing.T) {
	// Double all constants in each factor
	doubleConstants := func(expr Expression) Expression {
		if c, ok := expr.(*Constant); ok {
			return NewConstant(c.Value.Value * 2)
		}
		return expr
	}

	// Replace variable x with variable y
	replaceXWithY := func(expr Expression) Expression {
		if v, ok := expr.(*Variable); ok && v.Name == "x" {
			return NewVariable("y")
		}
		return expr
	}

	tests := []struct {
		name     string
		expr     Expression
		fn       func(Expression) Expression
		expected Expression
	}{
		{
			name:     "double single constant",
			expr:     NewConstant(5),
			fn:       doubleConstants,
			expected: NewConstant(10),
		},
		{
			name: "double constants in multiplication",
			expr: NewMultiply(
				NewConstant(2),
				NewConstant(3),
			),
			fn: doubleConstants,
			expected: NewMultiply(
				NewConstant(4),
				NewConstant(6),
			),
		},
		{
			name: "double constants in nested multiplication",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			),
			fn: doubleConstants,
			expected: NewMultiply(
				NewMultiply(
					NewConstant(4),
					NewConstant(6),
				),
				NewConstant(8),
			),
		},
		{
			name: "replace variable x with y",
			expr: NewMultiply(
				NewConstant(2),
				NewVariable("x"),
			),
			fn: replaceXWithY,
			expected: NewMultiply(
				NewConstant(2),
				NewVariable("y"),
			),
		},
		{
			name: "replace multiple x variables",
			expr: NewMultiply(
				NewMultiply(
					NewVariable("x"),
					NewConstant(3),
				),
				NewVariable("x"),
			),
			fn: replaceXWithY,
			expected: NewMultiply(
				NewMultiply(
					NewVariable("y"),
					NewConstant(3),
				),
				NewVariable("y"),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapFactors(tt.expr, tt.fn)

			if !expressionsEqual(result, tt.expected) {
				t.Errorf("MapFactors() result mismatch")
			}
		})
	}
}

// TestGetCoefficient tests extracting the coefficient (numeric part) of a monomial
func TestGetCoefficient(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected float64
		hasCoeff bool
	}{
		{
			name:     "constant only: 5",
			expr:     NewConstant(5),
			expected: 5,
			hasCoeff: true,
		},
		{
			name:     "variable only: x",
			expr:     NewVariable("x"),
			expected: 1,
			hasCoeff: true,
		},
		{
			name: "coefficient with variable: 3 * x",
			expr: NewMultiply(
				NewConstant(3),
				NewVariable("x"),
			),
			expected: 3,
			hasCoeff: true,
		},
		{
			name: "multiple constants: 2 * 3",
			expr: NewMultiply(
				NewConstant(2),
				NewConstant(3),
			),
			expected: 6,
			hasCoeff: true,
		},
		{
			name: "coefficient with multiple variables: 5 * x * y",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(5),
					NewVariable("x"),
				),
				NewVariable("y"),
			),
			expected: 5,
			hasCoeff: true,
		},
		{
			name: "variables only: x * y",
			expr: NewMultiply(
				NewVariable("x"),
				NewVariable("y"),
			),
			expected: 1,
			hasCoeff: true,
		},
		{
			name: "nested constants: (2 * 3) * x",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewConstant(3),
				),
				NewVariable("x"),
			),
			expected: 6,
			hasCoeff: true,
		},
		{
			name: "power expression: x^2",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(2),
			),
			expected: 1,
			hasCoeff: true,
		},
		{
			name: "coefficient with power: 4 * x^3",
			expr: NewMultiply(
				NewConstant(4),
				NewPower(
					NewVariable("x"),
					NewConstant(3),
				),
			),
			expected: 4,
			hasCoeff: true,
		},
		{
			name: "sqrt expression: sqrt(x)",
			expr: NewSqrt(
				NewVariable("x"),
			),
			expected: 1,
			hasCoeff: false,
		},
		{
			name: "coefficient with sqrt: 3 * sqrt(x)",
			expr: NewMultiply(
				NewConstant(3),
				NewSqrt(
					NewVariable("x"),
				),
			),
			expected: 3,
			hasCoeff: false,
		},
		{
			name: "complex coefficient: 2 * 3 * x (coefficient 6)",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewConstant(3),
				),
				NewVariable("x"),
			),
			expected: 6,
			hasCoeff: true,
		},
		{
			name: "coefficient with multiple powers: 5 * x^2 * y^3",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(5),
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
				),
				NewPower(
					NewVariable("y"),
					NewConstant(3),
				),
			),
			expected: 5,
			hasCoeff: true,
		},
		{
			name: "coefficient with power and sqrt: 2 * x^2 * sqrt(y)",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
				),
				NewSqrt(
					NewVariable("y"),
				),
			),
			expected: 2,
			hasCoeff: false,
		},
		{
			name: "fractional coefficient: 0.5 * x",
			expr: NewMultiply(
				NewConstant(0.5),
				NewVariable("x"),
			),
			expected: 0.5,
			hasCoeff: true,
		},
		{
			name: "negative coefficient: -3 * x",
			expr: NewMultiply(
				NewConstant(-3),
				NewVariable("x"),
			),
			expected: -3,
			hasCoeff: true,
		},
		{
			name: "nested power coefficient: 2 * (x^2)^3",
			expr: NewMultiply(
				NewConstant(2),
				NewPower(
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
					NewConstant(3),
				),
			),
			expected: 2,
			hasCoeff: true,
		},
		{
			name: "power with zero exponent: 3 * x^0 (mathematically 3 * 1 = 3)",
			expr: NewMultiply(
				NewConstant(3),
				NewPower(
					NewVariable("x"),
					NewConstant(0),
				),
			),
			expected: 3,
			hasCoeff: false,
		},
		{
			name: "nth root coefficient: 4 * cbrt(x)",
			expr: NewMultiply(
				NewConstant(4),
				NewNthRoot(
					NewVariable("x"),
					3,
				),
			),
			expected: 4,
			hasCoeff: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := GetCoefficient(tt.expr)

			if ok != tt.hasCoeff {
				t.Errorf("GetCoefficient() ok = %v, expected %v", ok, tt.hasCoeff)
				return
			}

			if ok && result != tt.expected {
				t.Errorf("GetCoefficient() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestGetDegree tests calculating the degree of a monomial (sum of variable exponents)
func TestGetDegree(t *testing.T) {
	tests := []struct {
		name      string
		expr      Expression
		expected  int
		hasDegree bool
	}{
		{
			name:      "constant only: 5",
			expr:      NewConstant(5),
			expected:  0,
			hasDegree: true,
		},
		{
			name:      "single variable: x",
			expr:      NewVariable("x"),
			expected:  1,
			hasDegree: true,
		},
		{
			name: "coefficient with variable: 3 * x",
			expr: NewMultiply(
				NewConstant(3),
				NewVariable("x"),
			),
			expected:  1,
			hasDegree: true,
		},
		{
			name: "two variables: x * y",
			expr: NewMultiply(
				NewVariable("x"),
				NewVariable("y"),
			),
			expected:  2,
			hasDegree: true,
		},
		{
			name: "coefficient with two variables: 2 * x * y",
			expr: NewMultiply(
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewVariable("y"),
			),
			expected:  2,
			hasDegree: true,
		},
		{
			name: "three variables: x * y * z",
			expr: NewMultiply(
				NewMultiply(
					NewVariable("x"),
					NewVariable("y"),
				),
				NewVariable("z"),
			),
			expected:  3,
			hasDegree: true,
		},
		{
			name: "coefficient with three variables: 5 * x * y * z",
			expr: NewMultiply(
				NewMultiply(
					NewMultiply(
						NewConstant(5),
						NewVariable("x"),
					),
					NewVariable("y"),
				),
				NewVariable("z"),
			),
			expected:  3,
			hasDegree: true,
		},
		{
			name: "only constants: 2 * 3",
			expr: NewMultiply(
				NewConstant(2),
				NewConstant(3),
			),
			expected:  0,
			hasDegree: true,
		},
		{
			name: "power expression: x^2 (mathematically degree 2)",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(2),
			),
			expected:  2,
			hasDegree: true,
		},
		{
			name: "power expression: x^3 (mathematically degree 3)",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(3),
			),
			expected:  3,
			hasDegree: true,
		},
		{
			name: "coefficient with power: 2 * x^3 (degree 3)",
			expr: NewMultiply(
				NewConstant(2),
				NewPower(
					NewVariable("x"),
					NewConstant(3),
				),
			),
			expected:  3,
			hasDegree: true,
		},
		{
			name: "two powers: x^2 * y^3 (mathematically degree 2+3=5)",
			expr: NewMultiply(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewPower(
					NewVariable("y"),
					NewConstant(3),
				),
			),
			expected:  5,
			hasDegree: true,
		},
		{
			name: "power and variable: x^2 * y (mathematically degree 2+1=3)",
			expr: NewMultiply(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewVariable("y"),
			),
			expected:  3,
			hasDegree: true,
		},
		{
			name: "complex: 3 * x^2 * y^3 * z (mathematically degree 2+3+1=6)",
			expr: NewMultiply(
				NewMultiply(
					NewMultiply(
						NewConstant(3),
						NewPower(
							NewVariable("x"),
							NewConstant(2),
						),
					),
					NewPower(
						NewVariable("y"),
						NewConstant(3),
					),
				),
				NewVariable("z"),
			),
			expected:  6,
			hasDegree: true,
		},
		{
			name: "power with zero exponent: x^0 (degree 0)",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(0),
			),
			expected:  0,
			hasDegree: true,
		},
		{
			name: "power with one exponent: x^1 (degree 1)",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(1),
			),
			expected:  1,
			hasDegree: true,
		},
		{
			name: "nested power: (x^2)^3 (mathematically x^6, degree 6)",
			expr: NewPower(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected:  6,
			hasDegree: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := GetDegree(tt.expr)

			if ok != tt.hasDegree {
				t.Errorf("GetDegree() ok = %v, expected %v", ok, tt.hasDegree)
				return
			}

			if ok && result != tt.expected {
				t.Errorf("GetDegree() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestGetDegreeWithFractionalExponents tests GetDegree with fractional exponents (if supported)
// Mathematically: x^(1/2) has degree 1/2, x^(3/2) has degree 3/2, etc.
func TestGetDegreeWithFractionalExponents(t *testing.T) {
	tests := []struct {
		name      string
		expr      Expression
		expected  float64 // Changed to float64 for fractional degrees
		hasDegree bool
	}{
		{
			name: "sqrt as x^(1/2): degree 0.5",
			expr: NewSqrt(
				NewVariable("x"),
			),
			expected:  0.5,
			hasDegree: true,
		},
		{
			name: "cube root as x^(1/3): degree 1/3",
			expr: NewNthRoot(
				NewVariable("x"),
				3,
			),
			expected:  1.0 / 3.0,
			hasDegree: true,
		},
		{
			name: "fourth root as x^(1/4): degree 1/4",
			expr: NewNthRoot(
				NewVariable("x"),
				4,
			),
			expected:  0.25,
			hasDegree: true,
		},
		{
			name: "x^2 * sqrt(x) = x^(5/2): degree 2.5",
			expr: NewMultiply(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewSqrt(
					NewVariable("x"),
				),
			),
			expected:  2.5,
			hasDegree: true,
		},
		{
			name: "sqrt(x) * sqrt(x) = x: degree 1",
			expr: NewMultiply(
				NewSqrt(
					NewVariable("x"),
				),
				NewSqrt(
					NewVariable("x"),
				),
			),
			expected:  1.0,
			hasDegree: true,
		},
		{
			name: "sqrt(x) * cbrt(x) = x^(1/2 + 1/3) = x^(5/6): degree 5/6",
			expr: NewMultiply(
				NewSqrt(
					NewVariable("x"),
				),
				NewNthRoot(
					NewVariable("x"),
					3,
				),
			),
			expected:  5.0 / 6.0,
			hasDegree: true,
		},
		{
			name: "sqrt(x) * sqrt(y): total degree 0.5 + 0.5 = 1",
			expr: NewMultiply(
				NewSqrt(
					NewVariable("x"),
				),
				NewSqrt(
					NewVariable("y"),
				),
			),
			expected:  1.0,
			hasDegree: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test is for future implementation
			// Current GetDegree returns int, but mathematically should handle fractional degrees
			t.Skip("Skipping fractional degree test - requires implementation that handles float64 degrees")
		})
	}
}

// TestMonomialSimplification tests mathematical simplifications (if implemented)
// These tests verify mathematical identities
func TestMonomialSimplification(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected Expression
		comment  string
	}{
		{
			name: "x^0 simplifies to 1",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(0),
			),
			expected: NewConstant(1),
			comment:  "Any non-zero number to power 0 equals 1",
		},
		{
			name: "x^1 simplifies to x",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(1),
			),
			expected: NewVariable("x"),
			comment:  "Any number to power 1 equals itself",
		},
		{
			name: "(sqrt(x))^2 simplifies to x",
			expr: NewPower(
				NewSqrt(
					NewVariable("x"),
				),
				NewConstant(2),
			),
			expected: NewVariable("x"),
			comment:  "Squaring a square root cancels out",
		},
		{
			name: "sqrt(x^2) simplifies to |x| or x (for positive x)",
			expr: NewSqrt(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
			),
			expected: NewVariable("x"),
			comment:  "Square root of square (assuming positive domain)",
		},
		{
			name: "(x^2)^3 simplifies to x^6",
			expr: NewPower(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected: NewPower(
				NewVariable("x"),
				NewConstant(6),
			),
			comment: "Power of power: (x^a)^b = x^(a*b)",
		},
		{
			name: "1 * x simplifies to x",
			expr: NewMultiply(
				NewConstant(1),
				NewVariable("x"),
			),
			expected: NewVariable("x"),
			comment:  "Multiplicative identity",
		},
		{
			name: "0 * x simplifies to 0",
			expr: NewMultiply(
				NewConstant(0),
				NewVariable("x"),
			),
			expected: NewConstant(0),
			comment:  "Multiplication by zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This test is for future implementation
			t.Skip("Skipping simplification test - requires Simplify() function implementation")
		})
	}
}

// TestEvalWithPowerAndRoot tests evaluation of expressions with powers and roots
func TestEvalWithPowerAndRoot(t *testing.T) {
	tests := []struct {
		name     string
		expr     Expression
		expected float64
		delta    float64 // For floating point comparison
		ok       bool
	}{
		{
			name: "2^3 = 8",
			expr: NewPower(
				NewConstant(2),
				NewConstant(3),
			),
			expected: 8.0,
			delta:    0.0001,
			ok:       true,
		},
		{
			name: "sqrt(4) = 2",
			expr: NewSqrt(
				NewConstant(4),
			),
			expected: 2.0,
			delta:    0.0001,
			ok:       true,
		},
		{
			name: "cbrt(8) = 2",
			expr: NewNthRoot(
				NewConstant(8),
				3,
			),
			expected: 2.0,
			delta:    0.0001,
			ok:       true,
		},
		{
			name: "(sqrt(2))^2 ≈ 2",
			expr: NewPower(
				NewSqrt(
					NewConstant(2),
				),
				NewConstant(2),
			),
			expected: 2.0,
			delta:    0.0001,
			ok:       true,
		},
		{
			name: "sqrt(2^2) = 2",
			expr: NewSqrt(
				NewPower(
					NewConstant(2),
					NewConstant(2),
				),
			),
			expected: 2.0,
			delta:    0.0001,
			ok:       true,
		},
		{
			name: "(2^2)^3 = 64",
			expr: NewPower(
				NewPower(
					NewConstant(2),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected: 64.0,
			delta:    0.0001,
			ok:       true,
		},
		{
			name: "2 * 3^2 = 18",
			expr: NewMultiply(
				NewConstant(2),
				NewPower(
					NewConstant(3),
					NewConstant(2),
				),
			),
			expected: 18.0,
			delta:    0.0001,
			ok:       true,
		},
		{
			name: "sqrt(16) * cbrt(8) = 4 * 2 = 8",
			expr: NewMultiply(
				NewSqrt(
					NewConstant(16),
				),
				NewNthRoot(
					NewConstant(8),
					3,
				),
			),
			expected: 8.0,
			delta:    0.0001,
			ok:       true,
		},
		{
			name: "2^0 = 1",
			expr: NewPower(
				NewConstant(2),
				NewConstant(0),
			),
			expected: 1.0,
			delta:    0.0001,
			ok:       true,
		},
		{
			name: "2^(-2) = 0.25",
			expr: NewPower(
				NewConstant(2),
				NewConstant(-2),
			),
			expected: 0.25,
			delta:    0.0001,
			ok:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, ok := tt.expr.Eval()
			if ok != tt.ok {
				t.Errorf("Eval() ok = %v, expected %v", ok, tt.ok)
				return
			}
			if ok {
				if num, ok := result.(*NumberValue); ok {
					diff := num.Value - tt.expected
					if diff < 0 {
						diff = -diff
					}
					if diff > tt.delta {
						t.Errorf("Eval() = %v, expected %v (within %v)", num.Value, tt.expected, tt.delta)
					}
				} else {
					t.Errorf("Eval() did not return NumberValue")
				}
			}
		})
	}
}
