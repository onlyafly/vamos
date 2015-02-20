package ast

import (
	"fmt"
	"strconv"
	"strings"

	"vamos/lang/token"
)

////////// Slice of Nodes

// Nodes type represents an array of Nodes.
type Nodes []Node

func (ns Nodes) String() string {
	return strings.Join(nodesToStrings([]Node(ns)), "\n")
}

////////// Node

// Node represents a parsed lisp node.
type Node interface {
	fmt.Stringer

	Equals(Node) bool
	TypeName() string
	Loc() *token.Location
}

////////// AnnotatedNode

// AnnotatedNode represents a node with an annotation (^foo) before it.
type AnnotatedNode interface {
	Node
	Annotation() Node
	SetAnnotation(n Node)
}

func displayAnnotation(an AnnotatedNode, rawRepresentation string) string {
	if an.Annotation() != nil {
		return "^" + an.Annotation().String() + " " + rawRepresentation
	}

	return rawRepresentation
}

////////// Expressions and Declarations

// Expr is a node representing an expression
type Expr interface {
	Node
	isExpr() bool
}

// Decl is a node representing a declaration
// TODO unused!
type Decl interface {
	Node
	isDecl() bool
}

////////// Str

type Str struct {
	Value      string
	annotation Node
	Location   *token.Location
}

func NewStr(value string) *Str { return &Str{Value: value} }

func (s *Str) String() string       { return displayAnnotation(s, "\""+s.Value+"\"") }
func (s *Str) isExpr() bool         { return true }
func (s *Str) Annotation() Node     { return s.annotation }
func (s *Str) SetAnnotation(n Node) { s.annotation = n }
func (s *Str) Equals(n Node) bool   { return s.Value == asStr(n).Value }
func (s *Str) TypeName() string     { return "string" }
func (s *Str) Loc() *token.Location { return s.Location }

////////// Nil

type Nil struct {
	Location   *token.Location
	annotation Node
}

func (n *Nil) String() string         { return "nil" }
func (n *Nil) isExpr() bool           { return true }
func (n *Nil) Annotation() Node       { return n.annotation }
func (n *Nil) SetAnnotation(ann Node) { n.annotation = ann }
func (n *Nil) TypeName() string       { return "nil" }
func (n *Nil) Loc() *token.Location   { return n.Location }
func (n *Nil) Equals(other Node) bool {
	if _, ok := other.(*Nil); ok {
		return true
	}
	return false
}

////////// Symbol

type Symbol struct {
	Name       string
	annotation Node
	Location   *token.Location
}

func (s *Symbol) String() string       { return displayAnnotation(s, s.Name) }
func (s *Symbol) isExpr() bool         { return true }
func (s *Symbol) Annotation() Node     { return s.annotation }
func (s *Symbol) SetAnnotation(n Node) { s.annotation = n }
func (s *Symbol) Equals(n Node) bool   { return s.Name == asSymbol(n).Name }
func (s *Symbol) TypeName() string     { return "symbol" }
func (s *Symbol) Loc() *token.Location { return s.Location }

////////// CharNode

// FIX rename to Char
type CharNode struct {
	Value      rune
	annotation Node
	Location   *token.Location
}

func (cn *CharNode) String() string {
	var rep string
	switch cn.Value {
	case '\n':
		rep = "\\newline"
	default:
		rep = fmt.Sprintf("\\%c", cn.Value)
	}
	return displayAnnotation(cn, rep)
}
func (cn *CharNode) isExpr() bool         { return true }
func (cn *CharNode) Annotation() Node     { return cn.annotation }
func (cn *CharNode) SetAnnotation(n Node) { cn.annotation = n }
func (cn *CharNode) Equals(n Node) bool   { return cn.Value == asChar(n).Value }
func (cn *CharNode) TypeName() string     { return "char" }
func (cn *CharNode) Loc() *token.Location { return cn.Location }

////////// Number

type Number struct {
	Value      float64
	annotation Node
	Location   *token.Location
}

func (num *Number) String() string {
	rep := strconv.FormatFloat(
		num.Value,
		'f',
		-1,
		64)

	return displayAnnotation(num, rep)
}

func (num *Number) isExpr() bool         { return true }
func (num *Number) Annotation() Node     { return num.annotation }
func (num *Number) SetAnnotation(n Node) { num.annotation = n }
func (num *Number) Equals(n Node) bool   { return num.Value == asNumber(n).Value }
func (num *Number) TypeName() string     { return "number" }
func (num *Number) Loc() *token.Location { return num.Location }

////////// List

type List struct {
	Nodes      []Node
	annotation Node
	Location   *token.Location
}

func NewList(nodes []Node) *List {
	return &List{Nodes: nodes}
}

func (l *List) String() string {
	raw := "(" + strings.Join(nodesToStrings(l.Nodes), " ") + ")"
	return displayAnnotation(l, raw)
}

func (l *List) isExpr() bool         { return true }
func (l *List) Annotation() Node     { return l.annotation }
func (l *List) SetAnnotation(n Node) { l.annotation = n }
func (l *List) TypeName() string     { return "list" }
func (l *List) Loc() *token.Location { return l.Location }
func (l *List) Equals(n Node) bool {
	other := asList(n)

	// Compare lengths
	if len(l.Nodes) != len(other.Nodes) {
		return false
	}

	// Compare contents
	for i, v := range l.Nodes {
		if !v.Equals(other.Nodes[i]) {
			return false
		}
	}

	return true
}

////////// Helpers

func asStr(n Node) *Str {
	if result, ok := n.(*Str); ok {
		return result
	}
	return &Str{}
}
func asChar(n Node) *CharNode {
	if result, ok := n.(*CharNode); ok {
		return result
	}
	return &CharNode{}
}
func asSymbol(n Node) *Symbol {
	if result, ok := n.(*Symbol); ok {
		return result
	}
	return &Symbol{}
}
func asNumber(n Node) *Number {
	if result, ok := n.(*Number); ok {
		return result
	}
	return &Number{}
}
func asList(n Node) *List {
	if result, ok := n.(*List); ok {
		return result
	}
	return &List{}
}

func nodesToStrings(nodes []Node) []string {
	return nodesToStringsWithFunc(nodes, func(n Node) string { return n.String() })
}

func nodesToStringsWithFunc(nodes []Node, convert func(n Node) string) []string {
	strings := make([]string, len(nodes))
	for i, node := range nodes {
		strings[i] = convert(node)
	}
	return strings
}
