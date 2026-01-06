package main

import (
	"exprtree/expr"
	"exprtree/latex"
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
			if result.Value != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result.Value)
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
			if result.Value != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result.Value)
			}
		})
	}
}

func TestIntegration_Errors(t *testing.T) {
	tests := []string{
		"2 +",      // Incomplete expression
		"(2 + 3",   // Unmatched paren
		"2 + + 3",  // Invalid syntax
		"",         // Empty input
		"* 2",      // Missing left operand
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
	expression, err := latex.ParseLatex(input)
	if err != nil {
		t.Fatalf("ParseLatex failed: %v", err)
	}

	// Verify we can traverse the expression tree
	children := expression.Children()
	if len(children) != 2 {
		t.Errorf("expected 2 children for AddExpression, got %d", len(children))
	}

	// Verify we can evaluate the expression
	result, ok := expression.Eval()
	if !ok {
		t.Errorf("evaluation failed")
	}

	numResult, ok := result.(*expr.NumberValue)
	if !ok || numResult.Value != 5.0 {
		t.Errorf("expected result 5.0, got %f", numResult.Value)
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
	if result.Value != expected {
		t.Errorf("expected %f, got %f", expected, result.Value)
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
			if result.Value != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, result.Value)
			}
		})
	}
}
