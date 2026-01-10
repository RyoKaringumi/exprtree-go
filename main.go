package main

import (
	"exprtree/latex"
	"fmt"
)

func main() {
	// 等式のテスト
	expression, err := latex.ParseLatex("2 + 3 = 5")
	if err != nil {
		fmt.Printf("Parse error: %v\n", err)
		return
	}

	result, ok := expression.Eval()
	if !ok {
		fmt.Printf("Evaluation failed\n")
		return
	}

	fmt.Printf("2 + 3 = 5 の評価結果: %+v\n", result)
}
