package expr

func SplitToTerms(node Expression) []Expression {
	var terms []Expression
	switch e := node.(type) {
	case *Add:
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
		result = NewAdd(result, terms[i])
	}
	return result
}

func CountTerms(node Expression) int {
	switch e := node.(type) {
	case *Add:
		return CountTerms(e.Left) + CountTerms(e.Right)
	default:
		return 1
	}
}

func IsPolynomialTerm(node Expression) bool {
	return IsMonomial(node)
}

func MapTerms(node Expression, fn func(Expression) Expression) Expression {
	terms := SplitToTerms(node)
	for i, term := range terms {
		terms[i] = fn(term)
	}
	return CombineTerms(terms)
}
