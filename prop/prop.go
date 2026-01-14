package prop

import "exprtree/expr"

type Proposition interface {
	expr.Expr
	Equals(other any) bool
}
