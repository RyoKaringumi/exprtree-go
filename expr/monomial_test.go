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
			expr: NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewMultiplyExpression(
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
			expr: NewAddExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected: []Expression{
				NewAddExpression(
					NewConstant(2),
					NewConstant(3),
				),
			},
		},
		{
			name: "division only: 10 / 2",
			expr: NewDivideExpression(
				NewConstant(10),
				NewConstant(2),
			),
			expected: []Expression{
				NewDivideExpression(
					NewConstant(10),
					NewConstant(2),
				),
			},
		},
		{
			name: "subtraction only: 5 - 3",
			expr: NewSubtractExpression(
				NewConstant(5),
				NewConstant(3),
			),
			expected: []Expression{
				NewSubtractExpression(
					NewConstant(5),
					NewConstant(3),
				),
			},
		},
		{
			name: "multiplication with addition: (2 + 3) * 4",
			expr: NewMultiplyExpression(
				NewAddExpression(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			),
			expected: []Expression{
				NewAddExpression(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			},
		},
		{
			name: "power expression: x^2",
			expr: NewPowerExpression(
				NewVariable("x"),
				NewConstant(2),
			),
			expected: []Expression{
				NewPowerExpression(
					NewVariable("x"),
					NewConstant(2),
				),
			},
		},
		{
			name: "multiplication with power: 3 * x^2",
			expr: NewMultiplyExpression(
				NewConstant(3),
				NewPowerExpression(
					NewVariable("x"),
					NewConstant(2),
				),
			),
			expected: []Expression{
				NewConstant(3),
				NewPowerExpression(
					NewVariable("x"),
					NewConstant(2),
				),
			},
		},
		{
			name: "sqrt expression: sqrt(x)",
			expr: NewSqrtExpression(
				NewVariable("x"),
			),
			expected: []Expression{
				NewSqrtExpression(
					NewVariable("x"),
				),
			},
		},
		{
			name: "multiplication with sqrt: 2 * sqrt(x) * y",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewConstant(2),
					NewSqrtExpression(
						NewVariable("x"),
					),
				),
				NewVariable("y"),
			),
			expected: []Expression{
				NewConstant(2),
				NewSqrtExpression(
					NewVariable("x"),
				),
				NewVariable("y"),
			},
		},
		{
			name: "nth root expression: cbrt(8)",
			expr: NewNthRootExpression(
				NewConstant(8),
				3,
			),
			expected: []Expression{
				NewNthRootExpression(
					NewConstant(8),
					3,
				),
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
			expected: NewMultiplyExpression(
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
			expected: NewMultiplyExpression(
				NewMultiplyExpression(
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
			expected: NewMultiplyExpression(
				NewMultiplyExpression(
					NewMultiplyExpression(
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
			expected: NewMultiplyExpression(
				NewMultiplyExpression(
					NewConstant(2),
					NewVariable("x"),
				),
				NewVariable("y"),
			),
		},
		{
			name: "complex factors with addition: (2 + 3) * x",
			factors: []Expression{
				NewAddExpression(
					NewConstant(2),
					NewConstant(3),
				),
				NewVariable("x"),
			},
			expected: NewMultiplyExpression(
				NewAddExpression(
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
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewConstant(3),
			),
		},
		{
			name: "nested multiplication: (2 * 3) * 4",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			),
		},
		{
			name: "variables and constants: 2 * x * y",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected: 2,
		},
		{
			name: "nested multiplication: (2 * 3) * 4",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			),
			expected: 3,
		},
		{
			name: "deeply nested: ((2 * 3) * 4) * 5",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewMultiplyExpression(
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
			expr: NewAddExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected: 1,
		},
		{
			name: "variables: x * y * z",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewVariable("x"),
			),
			expected: true,
		},
		{
			name: "nested multiplication: 2 * x * y",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewConstant(2),
					NewVariable("x"),
				),
				NewVariable("y"),
			),
			expected: true,
		},
		{
			name: "addition (not a monomial)",
			expr: NewAddExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected: false,
		},
		{
			name: "subtraction (not a monomial)",
			expr: NewSubtractExpression(
				NewConstant(5),
				NewConstant(3),
			),
			expected: false,
		},
		{
			name: "multiplication containing addition (not a monomial)",
			expr: NewMultiplyExpression(
				NewAddExpression(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected: false,
		},
		{
			name: "division with constant denominator",
			expr: NewDivideExpression(
				NewVariable("x"),
				NewConstant(2),
			),
			expected: true,
		},
		{
			name: "division with variable denominator (not a monomial in strict sense)",
			expr: NewDivideExpression(
				NewConstant(2),
				NewVariable("x"),
			),
			expected: false,
		},
		{
			name: "complex monomial: 3 * x * y * z",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewMultiplyExpression(
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
			expr: NewPowerExpression(
				NewVariable("x"),
				NewConstant(2),
			),
			expected: true,
		},
		{
			name: "power with coefficient: 3 * x^2",
			expr: NewMultiplyExpression(
				NewConstant(3),
				NewPowerExpression(
					NewVariable("x"),
					NewConstant(2),
				),
			),
			expected: true,
		},
		{
			name: "square root: sqrt(x)",
			expr: NewSqrtExpression(
				NewVariable("x"),
			),
			expected: true,
		},
		{
			name: "nth root: cbrt(x)",
			expr: NewNthRootExpression(
				NewVariable("x"),
				3,
			),
			expected: true,
		},
		{
			name: "coefficient with sqrt: 2 * sqrt(x)",
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewSqrtExpression(
					NewVariable("x"),
				),
			),
			expected: true,
		},
		{
			name: "power and variable: x^2 * y",
			expr: NewMultiplyExpression(
				NewPowerExpression(
					NewVariable("x"),
					NewConstant(2),
				),
				NewVariable("y"),
			),
			expected: true,
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
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewConstant(3),
			),
			fn: doubleConstants,
			expected: NewMultiplyExpression(
				NewConstant(4),
				NewConstant(6),
			),
		},
		{
			name: "double constants in nested multiplication",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewConstant(2),
					NewConstant(3),
				),
				NewConstant(4),
			),
			fn: doubleConstants,
			expected: NewMultiplyExpression(
				NewMultiplyExpression(
					NewConstant(4),
					NewConstant(6),
				),
				NewConstant(8),
			),
		},
		{
			name: "replace variable x with y",
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewVariable("x"),
			),
			fn: replaceXWithY,
			expected: NewMultiplyExpression(
				NewConstant(2),
				NewVariable("y"),
			),
		},
		{
			name: "replace multiple x variables",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewVariable("x"),
					NewConstant(3),
				),
				NewVariable("x"),
			),
			fn: replaceXWithY,
			expected: NewMultiplyExpression(
				NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewConstant(3),
				NewVariable("x"),
			),
			expected: 3,
			hasCoeff: true,
		},
		{
			name: "multiple constants: 2 * 3",
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected: 6,
			hasCoeff: true,
		},
		{
			name: "coefficient with multiple variables: 5 * x * y",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewVariable("x"),
				NewVariable("y"),
			),
			expected: 1,
			hasCoeff: true,
		},
		{
			name: "nested constants: (2 * 3) * x",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
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
			expr: NewPowerExpression(
				NewVariable("x"),
				NewConstant(2),
			),
			expected: 1,
			hasCoeff: true,
		},
		{
			name: "coefficient with power: 4 * x^3",
			expr: NewMultiplyExpression(
				NewConstant(4),
				NewPowerExpression(
					NewVariable("x"),
					NewConstant(3),
				),
			),
			expected: 4,
			hasCoeff: true,
		},
		{
			name: "sqrt expression: sqrt(x)",
			expr: NewSqrtExpression(
				NewVariable("x"),
			),
			expected: 1,
			hasCoeff: true,
		},
		{
			name: "coefficient with sqrt: 3 * sqrt(x)",
			expr: NewMultiplyExpression(
				NewConstant(3),
				NewSqrtExpression(
					NewVariable("x"),
				),
			),
			expected: 3,
			hasCoeff: true,
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
		name     string
		expr     Expression
		expected int
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
			expr: NewMultiplyExpression(
				NewConstant(3),
				NewVariable("x"),
			),
			expected:  1,
			hasDegree: true,
		},
		{
			name: "two variables: x * y",
			expr: NewMultiplyExpression(
				NewVariable("x"),
				NewVariable("y"),
			),
			expected:  2,
			hasDegree: true,
		},
		{
			name: "coefficient with two variables: 2 * x * y",
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewMultiplyExpression(
					NewMultiplyExpression(
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
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected:  0,
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
