package expr

func SplitToTerms(node Expression) []Expression {
	var terms []Expression
	switch e := node.(type) {
	case *AddExpression:
		terms = append(terms, SplitToTerms(e.Left)...)
		terms = append(terms, SplitToTerms(e.Right)...)
	default:
		terms = append(terms, node)
	}
	return terms
}

func CombineTerms(terms []Expression) Expression {
	if len(terms) == 0 {
		return nil
	}
	result := terms[0]
	for i := 1; i < len(terms); i++ {
		result = NewAddExpression(result, terms[i])
	}
	return result
}

func CountTerms(node Expression) int {
	switch e := node.(type) {
	case *AddExpression:
		return CountTerms(e.Left) + CountTerms(e.Right)
	default:
		return 1
	}
}

func IsPolynomialTerm(node Expression) bool {
	switch e := node.(type) {
	case *AddExpression, *SubtractExpression:
		return false
	case *MultiplyExpression:
		return IsPolynomialTerm(e.Left) && IsPolynomialTerm(e.Right)
	case *DivideExpression:
		_, ok := e.Right.Eval()
		return IsPolynomialTerm(e.Left) && ok
	case *Constant, *Variable:
		return true
	case *PowerExpression:
		return true
	case *SqrtExpression:
		return true
	default:
		return false
	}
}

func MapTerms(node Expression, fn func(Expression) Expression) Expression {
	terms := SplitToTerms(node)
	for i, term := range terms {
		terms[i] = fn(term)
	}
	return CombineTerms(terms)
}
