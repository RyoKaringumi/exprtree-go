package expr

import (
	"exprtree/ast"
	"exprtree/value"
)

type Expr interface {
	ast.HasChildren
	Eval() (value.Value, bool)
	Equals(other any) bool
}
