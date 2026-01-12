package expr

func SplitToFactors(node Expression) []Expression {
	var factors []Expression
	switch e := node.(type) {
	case *Multiply:
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
		result = NewMultiply(result, factors[i])
	}
	return result
}

func CountFactors(node Expression) int {
	switch e := node.(type) {
	case *Multiply:
		return CountFactors(e.Left) + CountFactors(e.Right)
	default:
		return 1
	}
}

func IsMonomial(node Expression) bool {
	switch e := node.(type) {
	case *Constant, *Variable:
		return true
	case *Multiply:
		return IsMonomial(e.Left) && IsMonomial(e.Right)
	case *Divide:
		_, ok := e.Right.Eval()
		return IsMonomial(e.Left) && ok

	case *Power:
		// べき乗の場合、指数が1以上の整数である必要がある
		// なぜならば、x^2の場合、xxに分解でき、結果として掛け算のみで構築されている事に出来る。
		// しかし、x^0.5などの場合、分解すると平方根が含まれてしまい、単項式ではなくなってしまう。
		exponentValue, ok := e.Right.Eval()
		if !ok {
			return false
		}
		if number, ok := exponentValue.(*NumberValue); ok {
			if number.Value >= 1 && number.Value == float64(int(number.Value)) {
				// 指数が1以上の整数の場合、基数が単項式である必要がある
				// さらに、基数が単項式である必要がある
				//
				// 基数がいかなる単項式だったとしても、整数のべき乗でその式を一定回数繰り返したものを展開するだけなので、単項式である
				return IsMonomial(e.Left)
			}
		}
		return false
	case *Sqrt:
		return false

	// 加算減算は多項式であって、単項式ではない
	case *Add, *Subtract:
		return false
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
	if !IsMonomial(node) {
		return 1, false
	}
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
		if powerExpr, ok := factor.(*Power); ok {
			exponentValue, ok := powerExpr.Right.Eval()
			if !ok {
				return 0, false
			}
			if number, ok := exponentValue.(*NumberValue); ok {
				if number.Value >= 0 && number.Value == float64(int(number.Value)) {
					leftDegree, ok := GetDegree(powerExpr.Left)
					if !ok {
						return 0, false
					}
					degree += int(number.Value) * leftDegree
				} else {
					return 0, false
				}
			} else {
				return 0, false
			}
		}
	}
	return degree, true
}
