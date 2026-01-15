package polynomial

import (
	"exprtree/expr"
	"exprtree/value"
)

func SplitMonomial(e expr.Expr) []expr.Expr {
	mul, ok := e.(*expr.Mul)
	if !ok {
		return []expr.Expr{e}
	}
	var factors []expr.Expr
	for _, child := range mul.Children() {
		expr, ok := child.(expr.Expr)
		if !ok {
			continue
		}
		factors = append(factors, SplitMonomial(expr)...)
	}
	return factors
}

func CombineMonomial(factors []expr.Expr) expr.Expr {
	if len(factors) == 0 {
		return nil
	}
	if len(factors) == 1 {
		return factors[0]
	}
	result := factors[0]
	for i := 1; i < len(factors); i++ {
		result = expr.NewMul(result, factors[i])
	}
	return result
}

func IsMonomial(e expr.Expr) bool {
	switch v := e.(type) {
	case *expr.Mul:
		return IsMonomial(v.Left()) && IsMonomial(v.Right())
	case *expr.Add:
		return false
	case *expr.Sub:
		return false
	case *expr.NthRoot:
		return false
	case *expr.Power:
		constant, ok := v.Exponent().(*expr.Constant)
		if !ok {
			return false
		}
		return IsMonomial(v.Base()) && value.IsPositiveIntegerReal(constant.Value())
	case *expr.Variable:
		return true
	case *expr.Constant:
		_, ok := v.Value().(*value.RealValue)
		if !ok {
			return false
		}
		return true
	default:
		return false
	}
}
