package main

import (
	"exprtree/expr"
	"exprtree/latex"
	"exprtree/value"
	"testing"
)

func TestIntegration_SimpleArithmetic(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"2 + 3", 5.0},
		{"10 - 4", 6.0},
		{"3 * 7", 21.0},
		{"15 / 3", 5.0},
		{"2 + 3 * 4", 14.0},
		{"(2 + 3) * 4", 20.0},
		{"10 - 2 - 3", 5.0},
		{"2.5 + 1.5", 4.0},
		{"(1 + 2) * (3 + 4)", 21.0},
		{"100 / 4 / 5", 5.0},
		{"2 * 3 + 4 * 5", 26.0},
		{"(10 + 20) / (3 + 2)", 6.0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := latex.ParseAndEval(tt.input)
			if err != nil {
				t.Fatalf("ParseAndEval failed: %v", err)
			}
			if result.Float64() != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result.Float64())
			}
		})
	}
}

func TestIntegration_DecimalNumbers(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14 + 2.86", 6.0},
		{"10.5 - 0.5", 10.0},
		{"2.5 * 4", 10.0},
		{"7.5 / 2.5", 3.0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := latex.ParseAndEval(tt.input)
			if err != nil {
				t.Fatalf("ParseAndEval failed: %v", err)
			}
			if result.Float64() != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result.Float64())
			}
		})
	}
}

func TestIntegration_Errors(t *testing.T) {
	tests := []string{
		"2 +",     // Incomplete expression
		"(2 + 3",  // Unmatched paren
		"2 + + 3", // Invalid syntax
		"",        // Empty input
		"* 2",     // Missing left operand
	}

	for _, input := range tests {
		t.Run(input, func(t *testing.T) {
			_, err := latex.ParseLatex(input)
			if err == nil {
				t.Errorf("expected error for input: %s", input)
			}
		})
	}
}

func TestIntegration_DivisionByZero(t *testing.T) {
	input := "10 / 0"
	result, err := latex.ParseAndEval(input)

	// Parsing should succeed
	if err != nil {
		// But evaluation should fail due to division by zero
		// This is acceptable - the error is caught during evaluation
		return
	}

	// If ParseAndEval succeeded, it means evaluation didn't detect division by zero
	// But our Eval() returns (nil, false) for division by zero, which ParseAndEval
	// should convert to an error
	if result != nil {
		t.Errorf("expected evaluation to fail for division by zero")
	}
}

func TestIntegration_ParseLatexReturnsExpression(t *testing.T) {
	input := "2 + 3"
	result, err := latex.ParseLatex(input)
	if err != nil {
		t.Fatalf("ParseLatex failed: %v", err)
	}

	// Type assert to expr.Expression
	expression, ok := result.(expr.Expr)
	if !ok {
		t.Fatalf("ParseLatex did not return an Expression, got %T", result)
	}

	// Verify we can evaluate the expression
	evalResult, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := evalResult.(*value.RealValue)
	if !ok || numResult.Float64() != 5.0 {
		t.Errorf("expected result 5.0, got %f", numResult.Float64())
	}
}

func TestIntegration_ComplexNesting(t *testing.T) {
	input := "((2 + 3) * (4 - 1)) / (7 - 4)"
	result, err := latex.ParseAndEval(input)
	if err != nil {
		t.Fatalf("ParseAndEval failed: %v", err)
	}

	// ((2 + 3) * (4 - 1)) / (7 - 4) = (5 * 3) / 3 = 15 / 3 = 5
	expected := 5.0
	if result.Float64() != expected {
		t.Errorf("expected %f, got %f", expected, result.Float64())
	}
}

func TestIntegration_Whitespace(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"  2  +  3  ", 5.0},
		{"\t10\t-\t4\t", 6.0},
		{"  ( 2 + 3 ) * 4  ", 20.0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := latex.ParseAndEval(tt.input)
			if err != nil {
				t.Fatalf("ParseAndEval failed: %v", err)
			}
			if result.Float64() != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result.Float64())
			}
		})
	}
}

func TestIntegration_Power(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"2^3", 8.0},
		{"2^3^2", 512.0},  // Right associative: 2^(3^2) = 2^9
		{"(2^3)^2", 64.0}, // Parentheses override: (2^3)^2 = 8^2
		{"2 + 3^4", 83.0}, // Precedence: 2 + (3^4) = 2 + 81
		{"2 * 3^2", 18.0}, // 2 * (3^2) = 2 * 9
		{"4^0.5", 2.0},    // Fractional exponent
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := latex.ParseAndEval(tt.input)
			if err != nil {
				t.Fatalf("Parse/eval error for '%s': %v", tt.input, err)
			}
			if result.Float64() != tt.expected {
				t.Errorf("For '%s': expected %f, got %f", tt.input, tt.expected, result.Float64())
			}
		})
	}
}

func TestIntegration_Sqrt(t *testing.T) {
	tests := []struct {
		input     string
		expected  float64
		tolerance float64
	}{
		{"\\sqrt{4}", 2.0, 0},
		{"\\sqrt{9}", 3.0, 0},
		{"\\sqrt{2.25}", 1.5, 0},
		{"\\sqrt[3]{8}", 2.0, 1e-10},
		{"\\sqrt[3]{27}", 3.0, 1e-10},
		{"\\sqrt[4]{16}", 2.0, 1e-10},
		{"2 + \\sqrt{9}", 5.0, 0},
		{"\\sqrt{9} * 2", 6.0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := latex.ParseAndEval(tt.input)
			if err != nil {
				t.Fatalf("Parse/eval error for '%s': %v", tt.input, err)
			}
			diff := result.Float64() - tt.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("For '%s': expected %f, got %f (diff: %e)", tt.input, tt.expected, result.Value, diff)
			}
		})
	}
}

func TestIntegration_PowerAndSqrt(t *testing.T) {
	tests := []struct {
		input     string
		expected  float64
		tolerance float64
	}{
		{"\\sqrt{3^2 + 4^2}", 5.0, 0},  // Pythagorean: sqrt(9+16) = sqrt(25)
		{"\\sqrt{2^4}", 4.0, 0},        // sqrt(16)
		{"2^\\sqrt{4}", 4.0, 0},        // 2^2
		{"(\\sqrt{9})^2", 9.0, 0},      // 3^2
		{"\\sqrt[3]{2^3}", 2.0, 1e-10}, // cbrt(8)
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := latex.ParseAndEval(tt.input)
			if err != nil {
				t.Fatalf("Parse/eval error for '%s': %v", tt.input, err)
			}
			diff := result.Float64() - tt.expected
			if diff < 0 {
				diff = -diff
			}
			if diff > tt.tolerance {
				t.Errorf("For '%s': expected %f, got %f (diff: %e)", tt.input, tt.expected, result.Value, diff)
			}
		})
	}
}
