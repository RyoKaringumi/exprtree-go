package expr

func SplitToFactors(node Expression) []Expression {
	var factors []Expression
	switch e := node.(type) {
	case *MultiplyExpression:
		factors = append(factors, SplitToFactors(e.Left)...)
		factors = append(factors, SplitToFactors(e.Right)...)
	default:
		factors = append(factors, node)
	}
	return factors
}

func CombineFactors(factors []Expression) Expression {
	if len(factors) == 0 {
		return nil
	}
	result := factors[0]
	for i := 1; i < len(factors); i++ {
		result = NewMultiplyExpression(result, factors[i])
	}
	return result
}

func CountFactors(node Expression) int {
	switch e := node.(type) {
	case *MultiplyExpression:
		return CountFactors(e.Left) + CountFactors(e.Right)
	default:
		return 1
	}
}

func IsMonomial(node Expression) bool {
	switch e := node.(type) {
	case *Constant, *Variable:
		return true
	case *MultiplyExpression:
		return IsMonomial(e.Left) && IsMonomial(e.Right)
	case *DivideExpression:
		_, ok := e.Right.Eval()
		return IsMonomial(e.Left) && ok
	case *PowerExpression:
		return true
	case *SqrtExpression:
		return true
	default:
		return false
	}
}

func MapFactors(node Expression, fn func(Expression) Expression) Expression {
	factors := SplitToFactors(node)
	for i, factor := range factors {
		factors[i] = fn(factor)
	}
	return CombineFactors(factors)
}

func GetCoefficient(node Expression) (float64, bool) {
	var factors = SplitToFactors(node)
	var coefficientFactors = []Expression{
		NewConstant(1),
	}
	for _, factor := range factors {
		if _, ok := factor.Eval(); ok {
			coefficientFactors = append(coefficientFactors, factor)
		}
	}
	var coefficient, ok = CombineFactors(coefficientFactors).Eval()
	if !ok {
		return 0, false
	}
	if number, ok := coefficient.(*NumberValue); ok {
		return number.Value, true
	}
	return 0, false
}

func GetDegree(node Expression) (int, bool) {
	var degree = 0
	var factors = SplitToFactors(node)
	for _, factor := range factors {
		if _, ok := factor.(*Variable); ok {
			degree += 1
		}
	}
	return degree, true
}
