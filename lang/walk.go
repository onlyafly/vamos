package lang

import "vamos/lang/ast"

type Visitor interface {
	Visit(node ast.Node) (childVisitor Visitor)
}

// DepthFirstWalk traverses the AST in depth-first order.
func DepthFirstWalk(visitor Visitor, node ast.Node) {
	childVisitor := visitor.Visit(node)

	if childVisitor != nil {
		var children ast.Nodes
		switch val := node.(type) {
		case ast.Coll:
			children = val.Children()
		}

		if len(children) > 0 {
			for _, child := range children {
				DepthFirstWalk(childVisitor, child)
			}
		}

		childVisitor.Visit(nil)
	}
}
