package lang

import ()

type Visitor interface {
	Visit(node Node) (childVisitor Visitor)
}

// Traverses the AST in depth-first order. 
func DepthFirstWalk(visitor Visitor, node Node) {
	childVisitor := visitor.Visit(node)

	if childVisitor != nil {
		children := node.Children()

		if len(children) > 0 {
			for _, child := range children {
				DepthFirstWalk(childVisitor, child)
			}
		}

		childVisitor.Visit(nil)
	}
}
