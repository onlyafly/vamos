package ast

import (
	"fmt"
	"strconv"
	"strings"
	. "vamos/lang/helpers"
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
	Children() []Node
	Equals(Node) bool
	TypeName() string
	Loc() *TokenLocation
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

////////// CollectionNode

type CollectionNode interface {
	Node
	Append(CollectionNode) CollectionNode
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

////////// StringNode

type StringNode struct {
	Value      string
	annotation Node
	Location   *TokenLocation
}

func NewStringNode(value string) *StringNode { return &StringNode{Value: value} }

func (s *StringNode) String() string       { return displayAnnotation(s, "\""+s.Value+"\"") }
func (s *StringNode) Children() []Node     { return nil }
func (s *StringNode) isExpr() bool         { return true }
func (s *StringNode) Annotation() Node     { return s.annotation }
func (s *StringNode) SetAnnotation(n Node) { s.annotation = n }
func (s *StringNode) Equals(n Node) bool   { return s.Value == asStringNode(n).Value }
func (s *StringNode) TypeName() string     { return "string" }
func (s *StringNode) Loc() *TokenLocation  { return s.Location }

////////// Symbol

type Symbol struct {
	Name       string
	annotation Node
	Location   *TokenLocation
}

func (s *Symbol) String() string       { return displayAnnotation(s, s.Name) }
func (s *Symbol) Children() []Node     { return nil }
func (s *Symbol) isExpr() bool         { return true }
func (s *Symbol) Annotation() Node     { return s.annotation }
func (s *Symbol) SetAnnotation(n Node) { s.annotation = n }
func (s *Symbol) Equals(n Node) bool   { return s.Name == asSymbol(n).Name }
func (s *Symbol) TypeName() string     { return "symbol" }
func (s *Symbol) Loc() *TokenLocation  { return s.Location }

////////// Number

type Number struct {
	Value      float64
	annotation Node
	Location   *TokenLocation
}

func (num *Number) String() string {
	rep := strconv.FormatFloat(
		num.Value,
		'f',
		-1,
		64)

	return displayAnnotation(num, rep)
}

func (num *Number) Children() []Node     { return nil }
func (num *Number) isExpr() bool         { return true }
func (num *Number) Annotation() Node     { return num.annotation }
func (num *Number) SetAnnotation(n Node) { num.annotation = n }
func (num *Number) Equals(n Node) bool   { return num.Value == asNumber(n).Value }
func (num *Number) TypeName() string     { return "number" }
func (num *Number) Loc() *TokenLocation  { return num.Location }

////////// List

type List struct {
	Nodes      []Node
	annotation Node
	Location   *TokenLocation
}

func NewList(nodes []Node) *List {
	return &List{Nodes: nodes}
}

func (l *List) String() string {
	raw := "(" + strings.Join(nodesToStrings(l.Nodes), " ") + ")"
	return displayAnnotation(l, raw)
}

func (l *List) Children() []Node     { return l.Nodes }
func (l *List) isExpr() bool         { return true }
func (l *List) Annotation() Node     { return l.annotation }
func (l *List) SetAnnotation(n Node) { l.annotation = n }
func (l *List) TypeName() string     { return "list" }
func (l *List) Loc() *TokenLocation  { return l.Location }
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

func asStringNode(n Node) *StringNode {
	if result, ok := n.(*StringNode); ok {
		return result
	}
	return &StringNode{}
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
