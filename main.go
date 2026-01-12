package main

import (
	"exprtree/expr"
	"exprtree/latex"
	"fmt"
)

func main() {
	fmt.Println("=== 単純な等式のテスト ===")
	testSimpleEquality()

	fmt.Println("\n=== 連続等号のテスト ===")
	testChainedEquality()

	fmt.Println("\n=== 4項連続等号のテスト ===")
	testFourTermEquality()
}

// testSimpleEquality tests a simple equality expression
func testSimpleEquality() {
	result, err := latex.ParseLatex("2 + 3 = 5")
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	// Simple equality returns an Equal (which is a Proposition)
	equalExpr, ok := result.(*expr.Equal)
	if !ok {
		fmt.Printf("Expected Equal, got %T\n", result)
		return
	}

	// Evaluate the equality
	evalResult, ok := equalExpr.Eval()
	if !ok {
		fmt.Printf("Evaluation failed\n")
		return
	}

	boolResult, ok := evalResult.(*expr.BoolValue)
	if !ok {
		fmt.Printf("Expected BoolValue, got %T\n", evalResult)
		return
	}

	fmt.Printf("入力: 2 + 3 = 5\n")
	fmt.Printf("構造: Equal(Add(2, 3), 5)\n")
	fmt.Printf("評価結果: %v\n", boolResult.Value)
}

// testChainedEquality tests a three-term chained equality
func testChainedEquality() {
	result, err := latex.ParseLatex("a = b = c")
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	// Chained equality returns an And proposition
	andExpr, ok := result.(*expr.And)
	if !ok {
		fmt.Printf("Expected And, got %T\n", result)
		return
	}

	fmt.Printf("入力: a = b = c\n")
	fmt.Printf("構造: And(Equal(a, b), Equal(b, c))\n")
	fmt.Printf("型: %T\n", andExpr)

	// Show the structure
	leftEqual, ok := andExpr.Left.(*expr.Equal)
	if ok {
		leftVar1, _ := leftEqual.Left.(*expr.Variable)
		leftVar2, _ := leftEqual.Right.(*expr.Variable)
		fmt.Printf("  左側: Equal(%s, %s)\n", leftVar1.Name, leftVar2.Name)
	}

	rightEqual, ok := andExpr.Right.(*expr.Equal)
	if ok {
		rightVar1, _ := rightEqual.Left.(*expr.Variable)
		rightVar2, _ := rightEqual.Right.(*expr.Variable)
		fmt.Printf("  右側: Equal(%s, %s)\n", rightVar1.Name, rightVar2.Name)
	}
}

// testFourTermEquality tests a four-term chained equality
func testFourTermEquality() {
	result, err := latex.ParseLatex("1 = 1 = 1 = 1")
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	// Four-term chained equality returns nested And propositions
	outerAnd, ok := result.(*expr.And)
	if !ok {
		fmt.Printf("Expected And, got %T\n", result)
		return
	}

	fmt.Printf("入力: 1 = 1 = 1 = 1\n")
	fmt.Printf("構造: And(And(Equal(1, 1), Equal(1, 1)), Equal(1, 1))\n")
	fmt.Printf("型: %T\n", outerAnd)

	// Show nested structure
	innerAnd, ok := outerAnd.Left.(*expr.And)
	if ok {
		fmt.Printf("  外側のAnd.Left: And (入れ子構造)\n")

		innerLeftEqual, ok := innerAnd.Left.(*expr.Equal)
		if ok {
			fmt.Printf("    内側のAnd.Left: Equal(1, 1)\n")
			_ = innerLeftEqual
		}

		innerRightEqual, ok := innerAnd.Right.(*expr.Equal)
		if ok {
			fmt.Printf("    内側のAnd.Right: Equal(1, 1)\n")
			_ = innerRightEqual
		}
	}

	rightEqual, ok := outerAnd.Right.(*expr.Equal)
	if ok {
		fmt.Printf("  外側のAnd.Right: Equal(1, 1)\n")
		_ = rightEqual
	}
}
