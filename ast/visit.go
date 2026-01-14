package ast

type HasChildren interface {
	Children() []HasChildren
}

func Walk(node HasChildren, fn func(HasChildren)) {
	fn(node)
	for _, child := range node.Children() {
		Walk(child, fn)
	}
}
