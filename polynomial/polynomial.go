package polynomial

import "exprtree/expr"

func SplitPolynomial(ex expr.Expr) []expr.Expr {
	var terms []expr.Expr
	switch e := ex.(type) {
	case *expr.Add:
		terms = append(terms, SplitPolynomial(e.Left())...)
		terms = append(terms, SplitPolynomial(e.Right())...)
	default:
		terms = append(terms, ex)
	}
	return terms
}

func CombinePolynomial(terms []expr.Expr) expr.Expr {
	if len(terms) == 0 {
		return nil
	}
	if len(terms) == 1 {
		return terms[0]
	}
	result := terms[0]
	for i := 1; i < len(terms); i++ {
		result = expr.NewAdd(result, terms[i])
	}
	return result
}

func IsPolynomial(e expr.Expr) bool {
	switch v := e.(type) {
	case *expr.Add:
		return IsPolynomial(v.Left()) && IsPolynomial(v.Right())
	case *expr.Sub:
		return IsPolynomial(v.Left()) && IsPolynomial(v.Right())
	default:
		return IsMonomial(v)
	}
}
