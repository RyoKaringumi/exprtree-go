package expr

type Binary interface {
	Expr
	Left() Expr
	Right() Expr
}
