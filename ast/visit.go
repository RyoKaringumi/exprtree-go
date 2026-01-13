package ast

type HasChildren interface {
	Children() []HasChildren
}
