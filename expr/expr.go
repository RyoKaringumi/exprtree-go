package expr

import "exprtree/value"

type Expr interface {
	Eval() (value.Value, bool)
}
