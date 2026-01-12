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
		// 単一の定数は既に項なので、それ自体だけが分解結果となる。
		{
			name:     "single constant",
			expr:     NewConstant(5),
			expected: []Expression{NewConstant(5)},
		},
		// 単一の変数は多項式の1項であるため、そのまま1つの項として扱われる。
		{
			name:     "single variable",
			expr:     NewVariable("x"),
			expected: []Expression{NewVariable("x")},
		},
		// 加算は項の和なので、2+3 は2つの項 (2 と 3) に分解される。
		{
			name: "simple addition: 2 + 3",
			expr: NewAdd(
				NewConstant(2),
				NewConstant(3),
			),
			expected: []Expression{
				NewConstant(2),
				NewConstant(3),
			},
		},
		// 加算のネストは分解すると平坦化できる： (1+2)+3 -> 1,2,3 の3項。
		{
			name: "nested addition: (1 + 2) + 3",
			expr: NewAdd(
				NewAdd(
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
		// 多重にネストされた加算も平坦化でき、各定数が独立した項となる。
		{
			name: "deeply nested: ((1 + 2) + 3) + 4",
			expr: NewAdd(
				NewAdd(
					NewAdd(
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
		// 各項が乗法式であっても加算の被演算子として独立した項として扱う。
		// したがって 2*x + 3*y は項 2*x と 3*y に分解される。
		{
			name: "addition with multiplication: 2*x + 3*y",
			expr: NewAdd(
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewMultiply(
					NewConstant(3),
					NewVariable("y"),
				),
			),
			expected: []Expression{
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewMultiply(
					NewConstant(3),
					NewVariable("y"),
				),
			},
		},
		// 左右の加算をそれぞれ展開し、全体を項列に平坦化する。
		// (x + 2*y) + (3 + z) -> x, 2*y, 3, z
		{
			name: "complex expression: (x + 2*y) + (3 + z)",
			expr: NewAdd(
				NewAdd(
					NewVariable("x"),
					NewMultiply(
						NewConstant(2),
						NewVariable("y"),
					),
				),
				NewAdd(
					NewConstant(3),
					NewVariable("z"),
				),
			),
			expected: []Expression{
				NewVariable("x"),
				NewMultiply(
					NewConstant(2),
					NewVariable("y"),
				),
				NewConstant(3),
				NewVariable("z"),
			},
		},
		// 加算が存在しない場合、式全体が単一の項と見なされる。
		// したがって乗法のみの 2*3 は1つの項として返される。
		{
			name: "multiplication only: 2 * 3",
			expr: NewMultiply(
				NewConstant(2),
				NewConstant(3),
			),
			expected: []Expression{
				NewMultiply(
					NewConstant(2),
					NewConstant(3),
				),
			},
		},
		// 除法のみも加算を含まないため単一項として扱われる。
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
		// 減算は2項の演算だが、ここでは加算同様に分解せず式全体を1項として扱う。
		// （別途分解ルールがあれば変わるが、現在の期待は単一項）
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
		// べき乗は項としてそのまま扱われる。x^2 + x は項 x^2 と項 x に分離される。
		{
			name: "power expression in addition: x^2 + x",
			expr: NewAdd(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewVariable("x"),
			),
			expected: []Expression{
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewVariable("x"),
			},
		},
		// 根号も式として独立した項になり得る。x + sqrt(x) は項 x と項 sqrt(x) に分解される。
		{
			name: "sqrt in polynomial: x + sqrt(x)",
			expr: NewAdd(
				NewVariable("x"),
				NewSqrt(
					NewVariable("x"),
				),
			),
			expected: []Expression{
				NewVariable("x"),
				NewSqrt(
					NewVariable("x"),
				),
			},
		},
		// 典型的な多項式の各項は係数とべき乗の積として表されるため、
		// 3*x^2 + 2*x + 1 は項 (3*x^2), (2*x), (1) に分解される。
		{
			name: "complex polynomial with power: 3*x^2 + 2*x + 1",
			expr: NewAdd(
				NewAdd(
					NewMultiply(
						NewConstant(3),
						NewPower(
							NewVariable("x"),
							NewConstant(2),
						),
					),
					NewMultiply(
						NewConstant(2),
						NewVariable("x"),
					),
				),
				NewConstant(1),
			),
			expected: []Expression{
				NewMultiply(
					NewConstant(3),
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
				),
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewConstant(1),
			},
		},
		{
			name: "nth root in addition: x + cbrt(x)",
			expr: NewAdd(
				NewVariable("x"),
				NewNthRoot(
					NewVariable("x"),
					3,
				),
			),
			expected: []Expression{
				NewVariable("x"),
				NewNthRoot(
					NewVariable("x"),
					3,
				),
			},
		},
		{
			name: "quadratic: x^2 + 2*x + 1",
			expr: NewAdd(
				NewAdd(
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
					NewMultiply(
						NewConstant(2),
						NewVariable("x"),
					),
				),
				NewConstant(1),
			),
			expected: []Expression{
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewConstant(1),
			},
		},
		{
			name: "cubic: x^3 + x^2 + x + 1",
			expr: NewAdd(
				NewAdd(
					NewAdd(
						NewPower(
							NewVariable("x"),
							NewConstant(3),
						),
						NewPower(
							NewVariable("x"),
							NewConstant(2),
						),
					),
					NewVariable("x"),
				),
				NewConstant(1),
			),
			expected: []Expression{
				NewPower(
					NewVariable("x"),
					NewConstant(3),
				),
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewVariable("x"),
				NewConstant(1),
			},
		},
		{
			name: "multivariate: x^2 + x*y + y^2",
			expr: NewAdd(
				NewAdd(
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
					NewMultiply(
						NewVariable("x"),
						NewVariable("y"),
					),
				),
				NewPower(
					NewVariable("y"),
					NewConstant(2),
				),
			),
			expected: []Expression{
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewMultiply(
					NewVariable("x"),
					NewVariable("y"),
				),
				NewPower(
					NewVariable("y"),
					NewConstant(2),
				),
			},
		},
		{
			name: "with multiple roots: sqrt(x) + sqrt(y) + sqrt(z)",
			expr: NewAdd(
				NewAdd(
					NewSqrt(
						NewVariable("x"),
					),
					NewSqrt(
						NewVariable("y"),
					),
				),
				NewSqrt(
					NewVariable("z"),
				),
			),
			expected: []Expression{
				NewSqrt(
					NewVariable("x"),
				),
				NewSqrt(
					NewVariable("y"),
				),
				NewSqrt(
					NewVariable("z"),
				),
			},
		},
		{
			name: "mixed powers and roots: x^3 + x^2 + sqrt(x) + 1",
			expr: NewAdd(
				NewAdd(
					NewAdd(
						NewPower(
							NewVariable("x"),
							NewConstant(3),
						),
						NewPower(
							NewVariable("x"),
							NewConstant(2),
						),
					),
					NewSqrt(
						NewVariable("x"),
					),
				),
				NewConstant(1),
			),
			expected: []Expression{
				NewPower(
					NewVariable("x"),
					NewConstant(3),
				),
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewSqrt(
					NewVariable("x"),
				),
				NewConstant(1),
			},
		},
		{
			name: "nested power in polynomial: (x^2)^3 + x",
			expr: NewAdd(
				NewPower(
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
					NewConstant(3),
				),
				NewVariable("x"),
			),
			expected: []Expression{
				NewPower(
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
					NewConstant(3),
				),
				NewVariable("x"),
			},
		},
		{
			name: "different nth roots: sqrt(x) + cbrt(x) + x^(1/4)",
			expr: NewAdd(
				NewAdd(
					NewSqrt(
						NewVariable("x"),
					),
					NewNthRoot(
						NewVariable("x"),
						3,
					),
				),
				NewNthRoot(
					NewVariable("x"),
					4,
				),
			),
			expected: []Expression{
				NewSqrt(
					NewVariable("x"),
				),
				NewNthRoot(
					NewVariable("x"),
					3,
				),
				NewNthRoot(
					NewVariable("x"),
					4,
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
			expected: NewAdd(
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
			expected: NewAdd(
				NewAdd(
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
			expected: NewAdd(
				NewAdd(
					NewAdd(
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
			expected: NewAdd(
				NewAdd(
					NewVariable("x"),
					NewVariable("y"),
				),
				NewConstant(5),
			),
		},
		{
			name: "complex terms: 2*x + 3*y",
			terms: []Expression{
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewMultiply(
					NewConstant(3),
					NewVariable("y"),
				),
			},
			expected: NewAdd(
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewMultiply(
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
			expr: NewAdd(
				NewConstant(2),
				NewConstant(3),
			),
		},
		{
			name: "nested addition: (1 + 2) + 3",
			expr: NewAdd(
				NewAdd(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
		},
		{
			name: "addition with multiplication: 2*x + 3*y",
			expr: NewAdd(
				NewMultiply(
					NewConstant(2),
					NewVariable("x"),
				),
				NewMultiply(
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
			expr: NewAdd(
				NewConstant(2),
				NewConstant(3),
			),
			expected: 2,
		},
		{
			name: "nested addition: (1 + 2) + 3",
			expr: NewAdd(
				NewAdd(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
			expected: 3,
		},
		{
			name: "deeply nested: ((1 + 2) + 3) + 4",
			expr: NewAdd(
				NewAdd(
					NewAdd(
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
			expr: NewMultiply(
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
			expr: NewMultiply(
				NewConstant(2),
				NewVariable("x"),
			),
			expected: true,
		},
		{
			name: "division",
			expr: NewDivide(
				NewConstant(10),
				NewConstant(2),
			),
			expected: true,
		},
		{
			name: "subtraction",
			expr: NewSubtract(
				NewConstant(5),
				NewConstant(3),
			),
			expected: false,
		},
		{
			name: "addition (not a term)",
			expr: NewAdd(
				NewConstant(2),
				NewConstant(3),
			),
			expected: false,
		},
		{
			name: "nested multiplication",
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
			name: "multiplication containing addition (not a term)",
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
			name: "sqrt expression: sqrt(x)",
			expr: NewSqrt(
				NewVariable("x"),
			),
			expected: false,
		},
		{
			name: "nth root: cbrt(8)",
			expr: NewNthRoot(
				NewConstant(8),
				3,
			),
			expected: false,
		},
		{
			name: "sqrt with coefficient: 2 * sqrt(x)",
			expr: NewMultiply(
				NewConstant(2),
				NewSqrt(
					NewVariable("x"),
				),
			),
			expected: false,
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
			name: "multiple powers: x^2 * y^3",
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
			expected: true,
		},
		{
			name: "power and sqrt: x^2 * sqrt(y)",
			expr: NewMultiply(
				NewPower(
					NewVariable("x"),
					NewConstant(2),
				),
				NewSqrt(
					NewVariable("y"),
				),
			),
			expected: false,
		},
		{
			name: "complex term: 2 * x^2 * y * sqrt(z)",
			expr: NewMultiply(
				NewMultiply(
					NewMultiply(
						NewConstant(2),
						NewPower(
							NewVariable("x"),
							NewConstant(2),
						),
					),
					NewVariable("y"),
				),
				NewSqrt(
					NewVariable("z"),
				),
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
			name: "negative power: x^(-2) (mathematically not a polynomial term)",
			expr: NewPower(
				NewVariable("x"),
				NewConstant(-2),
			),
			expected: false,
		},
		{
			name: "fractional power: x^(1/2) (mathematically not a polynomial term)",
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
			name: "power with variable exponent: x^y (not a polynomial term)",
			expr: NewPower(
				NewVariable("x"),
				NewVariable("y"),
			),
			expected: false,
		},
		{
			name: "sqrt of sum: sqrt(x + y) (not a polynomial term)",
			expr: NewSqrt(
				NewAdd(
					NewVariable("x"),
					NewVariable("y"),
				),
			),
			expected: false,
		},
		{
			name: "power of sum: (x + y)^2 (expands to polynomial but structure is not a term)",
			expr: NewPower(
				NewAdd(
					NewVariable("x"),
					NewVariable("y"),
				),
				NewConstant(2),
			),
			expected: false,
		},
		{
			name: "multiplication with addition: (x + y) * z (not a term)",
			expr: NewMultiply(
				NewAdd(
					NewVariable("x"),
					NewVariable("y"),
				),
				NewVariable("z"),
			),
			expected: false,
		},
		{
			name: "sqrt with addition inside: sqrt(x^2 + 1) (not a polynomial term)",
			expr: NewSqrt(
				NewAdd(
					NewPower(
						NewVariable("x"),
						NewConstant(2),
					),
					NewConstant(1),
				),
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
			expr: NewAdd(
				NewConstant(2),
				NewConstant(3),
			),
			fn: doubleConstants,
			expected: NewAdd(
				NewConstant(4),
				NewConstant(6),
			),
		},
		{
			name: "double constants in nested addition",
			expr: NewAdd(
				NewAdd(
					NewConstant(1),
					NewConstant(2),
				),
				NewConstant(3),
			),
			fn: doubleConstants,
			expected: NewAdd(
				NewAdd(
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
