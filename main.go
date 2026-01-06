package main

import (
	"exprtree/expr"
	"exprtree/latex"
	"fmt"
)

func main() {
	// Example 1: Manual expression tree construction
	println("Example 1: Manual expression tree construction")
	expression := expr.NewAddExpression(
		&expr.Constant{Value: expr.NumberValue{Value: 10}},
		expr.NewMultiplyExpression(
			&expr.Constant{Value: expr.NumberValue{Value: 2}},
			&expr.Constant{Value: expr.NumberValue{Value: 3}},
		),
	)

	result, ok := expression.Eval()
	if ok {
		if numResult, ok := result.(*expr.NumberValue); ok {
			println("Result:", numResult.Value) // Should print: Result: 16
		} else {
			println("Evaluation error")
		}
	} else {
		println("Evaluation failed")
	}

	// Example 2: LaTeX parser usage
	println("\nExample 2: LaTeX parser usage")

	// Simple arithmetic
	result2, err := latex.ParseAndEval("2 + 3 * 4")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("2 + 3 * 4 = %.2f\n", result2.Value) // Should print: 14.00
	}

	// With parentheses
	result3, err := latex.ParseAndEval("(2 + 3) * 4")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("(2 + 3) * 4 = %.2f\n", result3.Value) // Should print: 20.00
	}

	// Decimal numbers
	result4, err := latex.ParseAndEval("2.5 + 1.5")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("2.5 + 1.5 = %.2f\n", result4.Value) // Should print: 4.00
	}

	// Complex expression
	result5, err := latex.ParseAndEval("(1 + 2) * (3 + 4)")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("(1 + 2) * (3 + 4) = %.2f\n", result5.Value) // Should print: 21.00
	}

	// Example 3: Get expression tree for inspection
	println("\nExample 3: Expression tree inspection")
	expr2, err := latex.ParseLatex("2 + 3 * 4")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Expression parsed successfully\n")
		fmt.Printf("Number of children: %d\n", len(expr2.Children()))
		result6, ok := expr2.Eval()
		if ok {
			if numResult, ok := result6.(*expr.NumberValue); ok {
				fmt.Printf("Evaluation result: %.2f\n", numResult.Value)
			}
		}
	}

	// Example 4: Converting Expression tree to LaTeX string
	println("\nExample 4: Expression tree to LaTeX string")

	// Manual expression: (2 + 3) * 4
	manualExpr := expr.NewMultiplyExpression(
		expr.NewAddExpression(
			&expr.Constant{Value: expr.NumberValue{Value: 2}},
			&expr.Constant{Value: expr.NumberValue{Value: 3}},
		),
		&expr.Constant{Value: expr.NumberValue{Value: 4}},
	)

	latexStr, err := latex.ExpressionToLatex(manualExpr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Expression tree → LaTeX: %s\n", latexStr) // Should print: (2 + 3) * 4
	}

	// Round-trip test: Parse → Evaluate → Convert back to string
	println("\nExample 5: Round-trip conversion")
	originalStr := "10 - (3 - 2)"
	fmt.Printf("Original: %s\n", originalStr)

	parsedExpr, err := latex.ParseLatex(originalStr)
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
	} else {
		// Evaluate
		result7, ok := parsedExpr.Eval()
		if ok {
			if numResult, ok := result7.(*expr.NumberValue); ok {
				fmt.Printf("Evaluated: %.2f\n", numResult.Value)
			}
		}

		// Convert back to string
		reconstructed, err := latex.ExpressionToLatex(parsedExpr)
		if err != nil {
			fmt.Printf("Export error: %v\n", err)
		} else {
			fmt.Printf("Reconstructed: %s\n", reconstructed)
		}
	}

	variableExpr, err := latex.ParseLatex("x + 2")
	fmt.Printf("\nExample 6: Handling variables\n")
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
	} else {
		fmt.Printf("Parsed expression with variable successfully\n")
		result8, ok := variableExpr.Eval()
		if !ok {
			fmt.Printf("Evaluation failed due to variable\n")
		} else {
			if numResult, ok := result8.(*expr.NumberValue); ok {
				fmt.Printf("Evaluated: %.2f\n", numResult.Value)
			}
		}
	}
}
