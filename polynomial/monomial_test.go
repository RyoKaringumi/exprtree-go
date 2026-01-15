package polynomial_test

import (
	"exprtree/expr"
	"exprtree/latex"
	"exprtree/polynomial"
	"testing"
)

func TestIsMonomial(t *testing.T) {
	tests := []struct {
		expr     string
		expected bool
	}{
		{"3*x^2*y", true},
		{"5", true},
		{"x^3", true},
		{"2*x*y*z", true},
		{"x + y", false},
		{"x^{-2}", false},
		{"\\sqrt{x}", false},
		{"x^2 + 1", false},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			parsedData, err := latex.ParseLatex(tt.expr)
			parsedExpr, ok := parsedData.(expr.Expr)
			if !ok {
				t.Fatalf("Parsed data is not an exprtree.Expr for expression %s, type %T", tt.expr, parsedData)
			}
			if err != nil {
				t.Fatalf("Failed to parse expression %s: %v", tt.expr, err)
			}
			result := polynomial.IsMonomial(parsedExpr)
			if result != tt.expected {
				t.Errorf("IsMonomial(%s) = %v; want %v, type: %T", tt.expr, result, tt.expected, parsedExpr)
			}
		})
	}
}
