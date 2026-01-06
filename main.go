package main

import "fmt"

func main() {
	// Example 1: Manual expression tree construction
	println("Example 1: Manual expression tree construction")
	expr := NewAddExpression(
		&Constant{Value: NumberValue{Value: 10}},
		NewMultiplyExpression(
			&Constant{Value: NumberValue{Value: 2}},
			&Constant{Value: NumberValue{Value: 3}},
		),
	)

	result, ok := expr.Eval()
	if ok {
		if numResult, ok := result.(*NumberValue); ok {
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
	result2, err := ParseAndEval("2 + 3 * 4")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("2 + 3 * 4 = %.2f\n", result2.Value) // Should print: 14.00
	}

	// With parentheses
	result3, err := ParseAndEval("(2 + 3) * 4")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("(2 + 3) * 4 = %.2f\n", result3.Value) // Should print: 20.00
	}

	// Decimal numbers
	result4, err := ParseAndEval("2.5 + 1.5")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("2.5 + 1.5 = %.2f\n", result4.Value) // Should print: 4.00
	}

	// Complex expression
	result5, err := ParseAndEval("(1 + 2) * (3 + 4)")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("(1 + 2) * (3 + 4) = %.2f\n", result5.Value) // Should print: 21.00
	}

	// Example 3: Get expression tree for inspection
	println("\nExample 3: Expression tree inspection")
	expr2, err := ParseLatex("2 + 3 * 4")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Expression parsed successfully\n")
		fmt.Printf("Number of children: %d\n", len(expr2.Children()))
		result6, ok := expr2.Eval()
		if ok {
			if numResult, ok := result6.(*NumberValue); ok {
				fmt.Printf("Evaluation result: %.2f\n", numResult.Value)
			}
		}
	}
}
