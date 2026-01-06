package main

func main() {
	// Example usage:
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
}
