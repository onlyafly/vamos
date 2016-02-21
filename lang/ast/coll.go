package ast

import (
	"errors"
	"fmt"
	"unicode/utf8"
)

type Coll interface {
	Node

	Children() []Node
	Append(coll Coll) (Coll, error)
	Cons(elem Node) (Coll, error)
	First() Node
	Rest() Node
	Length() int
	IsEmpty() bool
}

////////// Length

func (l *List) Length() int { return len(l.Nodes) }
func (n *Nil) Length() int  { return 0 }
func (s *Str) Length() int  { return utf8.RuneCountInString(s.Value) }

////////// Children

func (l *List) Children() []Node { return l.Nodes }
func (n *Nil) Children() []Node  { return nil }
func (s *Str) Children() []Node {
	if len(s.Value) == 0 {
		return []Node{}
	}
	cs := make([]Node, 0)
	for _, r := range s.Value {
		c := &Char{Value: r}
		cs = append(cs, c)
	}
	return cs
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

func (l *List) Append(other Coll) (Coll, error) {
	if other.IsEmpty() {
		return l, nil
	} else {
		return NewList(append(l.Nodes, other.Children()...)), nil
	}
}
func (n *Nil) Append(other Coll) (Coll, error) {
	return other, nil
}
func (s *Str) Append(other Coll) (Coll, error) {
	if other.IsEmpty() {
		return s, nil
	}

	switch val := other.(type) {
	case *Str:
		return NewStr(s.Value + val.Value), nil
	default:
		return nil, errors.New("Cannot append a non-string onto a string: " + val.String())
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
