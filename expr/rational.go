package expr

func NewFractionExpression(numerator Expression, denominator Expression) *Divide {
	return NewDivide(numerator, denominator)
}

func IsFractional(node Expression) bool {
	if _, ok := node.(*Divide); ok {
		return true
	}
	return false
}

// 分数の加算
func AddFractions(a, b *Divide) *Divide {
	return NewFractionExpression(
		NewAdd(
			NewMultiply(a.Left, b.Right),
			NewMultiply(b.Left, a.Right),
		),
		NewMultiply(a.Right, b.Right),
	)
}

// 分数の乗算
func MultiplyFractions(a, b *Divide) *Divide {
	return NewFractionExpression(
		NewMultiply(a.Left, b.Left),
		NewMultiply(a.Right, b.Right),
	)
}