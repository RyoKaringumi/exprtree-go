package expr

import (
	"testing"
)

// TestSplitToTerms tests the SplitToTerms function that decomposes an expression into additive terms
func TestSplitToTerms(t *testing.T) {
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
			name: "simple addition: 2 + 3",
			expr: NewAddExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected: []Expression{
				NewConstant(2),
				NewConstant(3),
			},
		},
		{
			name: "nested addition: (1 + 2) + 3",
			expr: NewAddExpression(
				NewAddExpression(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected: []Expression{
				NewConstant(1),
				NewConstant(2),
				NewConstant(3),
			},
		},
		{
			name: "deeply nested: ((1 + 2) + 3) + 4",
			expr: NewAddExpression(
				NewAddExpression(
					NewAddExpression(
						NewConstant(1),
						NewConstant(2),
					),
					NewConstant(3),
				),
				NewConstant(4),
			),
			expected: []Expression{
				NewConstant(1),
				NewConstant(2),
				NewConstant(3),
				NewConstant(4),
			},
		},
		{
			name: "addition with multiplication: 2*x + 3*y",
			expr: NewAddExpression(
				NewMultiplyExpression(
					NewConstant(2),
					NewVariable("x"),
				),
				NewMultiplyExpression(
					NewConstant(3),
					NewVariable("y"),
				),
			),
			expected: []Expression{
				NewMultiplyExpression(
					NewConstant(2),
					NewVariable("x"),
				),
				NewMultiplyExpression(
					NewConstant(3),
					NewVariable("y"),
				),
			},
		},
		{
			name: "complex expression: (x + 2*y) + (3 + z)",
			expr: NewAddExpression(
				NewAddExpression(
					NewVariable("x"),
					NewMultiplyExpression(
						NewConstant(2),
						NewVariable("y"),
					),
				),
				NewAddExpression(
					NewConstant(3),
					NewVariable("z"),
				),
			),
			expected: []Expression{
				NewVariable("x"),
				NewMultiplyExpression(
					NewConstant(2),
					NewVariable("y"),
				),
				NewConstant(3),
				NewVariable("z"),
			},
		},
		{
			name: "multiplication only: 2 * 3",
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected: []Expression{
				NewMultiplyExpression(
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SplitToTerms(tt.expr)

			if len(result) != len(tt.expected) {
				t.Errorf("SplitToTerms() returned %d terms, expected %d", len(result), len(tt.expected))
				return
			}

			for i, term := range result {
				if !expressionsEqual(term, tt.expected[i]) {
					t.Errorf("SplitToTerms() term[%d] mismatch", i)
				}
			}
		})
	}
}

// TestCombineTerms tests the CombineTerms function that combines a slice of expressions using addition
func TestCombineTerms(t *testing.T) {
	tests := []struct {
		name     string
		terms    []Expression
		expected Expression
	}{
		{
			name:     "single term: 5",
			terms:    []Expression{NewConstant(5)},
			expected: NewConstant(5),
		},
		{
			name: "two terms: 2 + 3",
			terms: []Expression{
				NewConstant(2),
				NewConstant(3),
			},
			expected: NewAddExpression(
				NewConstant(2),
				NewConstant(3),
			),
		},
		{
			name: "three terms: 1 + 2 + 3",
			terms: []Expression{
				NewConstant(1),
				NewConstant(2),
				NewConstant(3),
			},
			expected: NewAddExpression(
				NewAddExpression(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
		},
		{
			name: "four terms: 1 + 2 + 3 + 4",
			terms: []Expression{
				NewConstant(1),
				NewConstant(2),
				NewConstant(3),
				NewConstant(4),
			},
			expected: NewAddExpression(
				NewAddExpression(
					NewAddExpression(
						NewConstant(1),
						NewConstant(2),
					),
					NewConstant(3),
				),
				NewConstant(4),
			),
		},
		{
			name: "variables and constants: x + y + 5",
			terms: []Expression{
				NewVariable("x"),
				NewVariable("y"),
				NewConstant(5),
			},
			expected: NewAddExpression(
				NewAddExpression(
					NewVariable("x"),
					NewVariable("y"),
				),
				NewConstant(5),
			),
		},
		{
			name: "complex terms: 2*x + 3*y",
			terms: []Expression{
				NewMultiplyExpression(
					NewConstant(2),
					NewVariable("x"),
				),
				NewMultiplyExpression(
					NewConstant(3),
					NewVariable("y"),
				),
			},
			expected: NewAddExpression(
				NewMultiplyExpression(
					NewConstant(2),
					NewVariable("x"),
				),
				NewMultiplyExpression(
					NewConstant(3),
					NewVariable("y"),
				),
			),
		},
		{
			name:     "empty slice",
			terms:    []Expression{},
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CombineTerms(tt.terms)

			// Handle nil case
			if tt.expected == nil {
				if result != nil {
					t.Errorf("CombineTerms() expected nil, got non-nil")
				}
				return
			}

			if result == nil {
				t.Errorf("CombineTerms() returned nil, expected non-nil")
				return
			}

			if !expressionsEqual(result, tt.expected) {
				t.Errorf("CombineTerms() result mismatch")
			}
		})
	}
}

// TestRoundTripSplitCombine tests that splitting and combining terms is an identity operation
func TestRoundTripSplitCombine(t *testing.T) {
	tests := []struct {
		name string
		expr Expression
	}{
		{
			name: "simple addition: 2 + 3",
			expr: NewAddExpression(
				NewConstant(2),
				NewConstant(3),
			),
		},
		{
			name: "nested addition: (1 + 2) + 3",
			expr: NewAddExpression(
				NewAddExpression(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
		},
		{
			name: "addition with multiplication: 2*x + 3*y",
			expr: NewAddExpression(
				NewMultiplyExpression(
					NewConstant(2),
					NewVariable("x"),
				),
				NewMultiplyExpression(
					NewConstant(3),
					NewVariable("y"),
				),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Split and then combine
			terms := SplitToTerms(tt.expr)
			result := CombineTerms(terms)

			// Check if the result is structurally equal to the original
			if !expressionsEqual(result, tt.expr) {
				t.Errorf("Round trip failed: split and combine did not produce equivalent expression")
			}
		})
	}
}

// TestCountTerms tests the CountTerms function that counts the number of additive terms
func TestCountTerms(t *testing.T) {
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
			name: "simple addition: 2 + 3",
			expr: NewAddExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected: 2,
		},
		{
			name: "nested addition: (1 + 2) + 3",
			expr: NewAddExpression(
				NewAddExpression(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected: 3,
		},
		{
			name: "deeply nested: ((1 + 2) + 3) + 4",
			expr: NewAddExpression(
				NewAddExpression(
					NewAddExpression(
						NewConstant(1),
						NewConstant(2),
					),
					NewConstant(3),
				),
				NewConstant(4),
			),
			expected: 4,
		},
		{
			name: "multiplication only",
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CountTerms(tt.expr)

			if result != tt.expected {
				t.Errorf("CountTerms() = %d, expected %d", result, tt.expected)
			}
		})
	}
}

// TestIsPolynomialTerm tests whether an expression is a valid polynomial term (no addition)
func TestIsPolynomialTerm(t *testing.T) {
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
			name: "multiplication",
			expr: NewMultiplyExpression(
				NewConstant(2),
				NewVariable("x"),
			),
			expected: true,
		},
		{
			name: "division",
			expr: NewDivideExpression(
				NewConstant(10),
				NewConstant(2),
			),
			expected: true,
		},
		{
			name: "subtraction",
			expr: NewSubtractExpression(
				NewConstant(5),
				NewConstant(3),
			),
			expected: false,
		},
		{
			name: "addition (not a term)",
			expr: NewAddExpression(
				NewConstant(2),
				NewConstant(3),
			),
			expected: false,
		},
		{
			name: "nested multiplication",
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
			name: "multiplication containing addition (not a term)",
			expr: NewMultiplyExpression(
				NewAddExpression(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsPolynomialTerm(tt.expr)

			if result != tt.expected {
				t.Errorf("IsPolynomialTerm() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

// TestMapTerms tests the MapTerms function that applies a transformation to each term
func TestMapTerms(t *testing.T) {
	// Double all constants in each term
	doubleConstants := func(expr Expression) Expression {
		if c, ok := expr.(*Constant); ok {
			return NewConstant(c.Value.Value * 2)
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
			name: "double constants in addition",
			expr: NewAddExpression(
				NewConstant(2),
				NewConstant(3),
			),
			fn: doubleConstants,
			expected: NewAddExpression(
				NewConstant(4),
				NewConstant(6),
			),
		},
		{
			name: "double constants in nested addition",
			expr: NewAddExpression(
				NewAddExpression(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
			fn: doubleConstants,
			expected: NewAddExpression(
				NewAddExpression(
					NewConstant(2),
					NewConstant(4),
				),
				NewConstant(6),
			),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := MapTerms(tt.expr, tt.fn)

			if !expressionsEqual(result, tt.expected) {
				t.Errorf("MapTerms() result mismatch")
			}
		})
	}
}
