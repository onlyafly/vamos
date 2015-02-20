package lang

import (
	"fmt"
	"unicode/utf8"
	"vamos/lang/ast"
)

type Coll interface {
	Children() []ast.Node
	Append(coll Coll) Coll
	Cons(elem ast.Node) Coll
	First() ast.Node
	Rest() ast.Node
	IsEmpty() bool
}

type (
	ListColl ast.List
	StrColl  ast.Str
	NilColl  ast.Nil
)

////////// Children

func (l *ListColl) Children() []ast.Node { return l.Nodes }
func (n *NilColl) Children() []ast.Node  { return nil }
func (s *StrColl) Children() []ast.Node  { return nil /*TODO*/ }

////////// First

func (l *ListColl) First() ast.Node {
	if len(l.Nodes) == 0 {
		return &ast.Nil{}
	}
	return l.Nodes[0]
}
func (n *NilColl) First() ast.Node { return (*ast.Nil)(n) }
func (s *StrColl) First() ast.Node {
	if len(s.Value) == 0 {
		return &ast.Nil{}
	}
	r, _ := utf8.DecodeRuneInString(s.Value)
	return &ast.CharNode{Value: r}
}

////////// Rest

func (l *ListColl) Rest() ast.Node {
	if len(l.Nodes) == 0 {
		return (*ast.List)(l)
	}
	return ast.NewList(l.Nodes[1:])
}
func (n *NilColl) Rest() ast.Node { return &ast.List{} }
func (s *StrColl) Rest() ast.Node {
	if len(s.Value) == 0 {
		return (*ast.Str)(s)
	}
	_, firstRuneWidth := utf8.DecodeRuneInString(s.Value)
	return ast.NewStr(s.Value[firstRuneWidth:])
}

////////// Append

func (l *ListColl) Append(other Coll) Coll {
	if other.IsEmpty() {
		return l
	} else {
		return (*ListColl)(ast.NewList(append(l.Nodes, other.Children()...)))
	}
}
func (n *NilColl) Append(other Coll) Coll {
	return other
}
func (s *StrColl) Append(other Coll) Coll {
	if other.IsEmpty() {
		return s
	}

	switch val := other.(type) {
	case *StrColl:
		return (*StrColl)(ast.NewStr(s.Value + val.Value))
	case ast.Node:
		panic("Unrecognized collection type: " + val.String())
	default:
		panic("Unrecognized object")
	}
}

////////// Cons

func (l *ListColl) Cons(elem ast.Node) Coll {
	return (*ListColl)(ast.NewList(append([]ast.Node{elem}, l.Nodes...)))
}
func (n *NilColl) Cons(elem ast.Node) Coll {
	return (*ListColl)(ast.NewList([]ast.Node{elem}))
}
func (s *StrColl) Cons(elem ast.Node) Coll {
	switch val := elem.(type) {
	case *ast.CharNode:
		return (*StrColl)(ast.NewStr(fmt.Sprintf("%c%v", val.Value, s.Value)))
	}

	panicEvalError(s, "Cannot cons a non-character onto a string: "+elem.String())
	return nil
}

////////// IsEmpty

func (l *ListColl) IsEmpty() bool {
	return len(l.Nodes) == 0
}
func (n *NilColl) IsEmpty() bool { return true }
func (s *StrColl) IsEmpty() bool { return len(s.Value) == 0 }
