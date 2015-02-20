package ast

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

type Coll interface {
	Node

	Children() []Node
	Append(coll Coll) Coll
	Cons(elem Node) (Coll, error)
	First() Node
	Rest() Node
	IsEmpty() bool
}

////////// Children

func (l *List) Children() []Node { return l.Nodes }
func (n *Nil) Children() []Node  { return nil }
func (s *Str) Children() []Node {
	return nil /*TODO*/
}

////////// First

func (l *List) First() Node {
	if len(l.Nodes) == 0 {
		return &Nil{}
	}
	return l.Nodes[0]
}
func (n *Nil) First() Node { return n }
func (s *Str) First() Node {
	if len(s.Value) == 0 {
		return &Nil{}
	}
	r, _ := utf8.DecodeRuneInString(s.Value)
	return &Char{Value: r}
}

////////// Rest

func (lc *List) Rest() Node {
	if len(lc.Nodes) == 0 {
		return lc
	}
	return NewList(lc.Nodes[1:])
}
func (n *Nil) Rest() Node { return &List{} }
func (sc *Str) Rest() Node {
	if len(sc.Value) == 0 {
		return sc
	}
	_, firstRuneWidth := utf8.DecodeRuneInString(sc.Value)
	return NewStr(sc.Value[firstRuneWidth:])
}

////////// Append

func (l *List) Append(other Coll) Coll {
	if other.IsEmpty() {
		return l
	} else {
		return NewList(append(l.Nodes, other.Children()...))
	}
}
func (n *Nil) Append(other Coll) Coll {
	return other
}
func (s *Str) Append(other Coll) Coll {
	if other.IsEmpty() {
		return s
	}

	switch val := other.(type) {
	case *Str:
		return NewStr(s.Value + val.Value)
	case Node:
		panic("Unrecognized collection type: " + val.String())
	default:
		panic("Unrecognized object")
	}
}

////////// Cons

func (l *List) Cons(elem Node) (Coll, error) {
	return NewList(append([]Node{elem}, l.Nodes...)), nil
}
func (n *Nil) Cons(elem Node) (Coll, error) {
	return NewList([]Node{elem}), nil
}
func (s *Str) Cons(elem Node) (Coll, error) {
	switch val := elem.(type) {
	case *Char:
		return NewStr(fmt.Sprintf("%c%v", val.Value, s.Value)), nil
	}
	return nil, errors.New("Cannot cons a non-character onto a string: " + elem.String())
}

////////// IsEmpty

func (l *List) IsEmpty() bool {
	return len(l.Nodes) == 0
}
func (n *Nil) IsEmpty() bool { return true }
func (s *Str) IsEmpty() bool { return len(s.Value) == 0 }
