package expr

func NewFractionExpression(numerator Expression, denominator Expression) *DivideExpression {
	return NewDivideExpression(numerator, denominator)
}

func IsFractional(node Expression) bool {
	if _, ok := node.(*DivideExpression); ok {
		return true
	}
	return false
}

// 分数の加算
func AddFractions(a, b *DivideExpression) *DivideExpression {
	return NewFractionExpression(
		NewAddExpression(
			NewMultiplyExpression(a.Left, b.Right),
			NewMultiplyExpression(b.Left, a.Right),
		),
		NewMultiplyExpression(a.Right, b.Right),
	)
}

// 分数の乗算
func MultiplyFractions(a, b *DivideExpression) *DivideExpression {
	return NewFractionExpression(
		NewMultiplyExpression(a.Left, b.Left),
		NewMultiplyExpression(a.Right, b.Right),
	)
}