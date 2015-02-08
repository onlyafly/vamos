package lang

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"
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

////////// Coll

type Coll interface {
	Node
	Append(coll Coll) Coll
	Cons(elem Node) Coll
	First() Node
	Rest() Node
	IsEmpty() bool
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
func (s *StringNode) IsEmpty() bool        { return len(s.Value) == 0 }
func (s *StringNode) First() Node {
	if len(s.Value) == 0 {
		return &NilNode{}
	}
	r, _ := utf8.DecodeRuneInString(s.Value)
	return &CharNode{Value: r}
}
func (s *StringNode) Rest() Node {
	if len(s.Value) == 0 {
		return s
	}
	_, firstRuneWidth := utf8.DecodeRuneInString(s.Value)
	return NewStringNode(s.Value[firstRuneWidth:])
}
func (s *StringNode) Append(other Coll) Coll {
	if other.IsEmpty() {
		return s
	}

	switch val := other.(type) {
	case *StringNode:
		return NewStringNode(s.Value + val.Value)
	default:
		panic("Unrecognized collection type: " + val.String())
	}
}
func (s *StringNode) Cons(elem Node) Coll {
	switch val := elem.(type) {
	case *CharNode:
		return NewStringNode(fmt.Sprintf("%c%v", val.Value, s.Value))
	}

	panicEvalError(s, "Cannot cons a non-character onto a string: "+elem.String())
	return nil
}

////////// NilNode

type NilNode struct {
	Location   *TokenLocation
	annotation Node
}

func (n *NilNode) String() string         { return "nil" }
func (n *NilNode) Children() []Node       { return nil }
func (n *NilNode) isExpr() bool           { return true }
func (n *NilNode) Annotation() Node       { return n.annotation }
func (n *NilNode) SetAnnotation(ann Node) { n.annotation = ann }
func (n *NilNode) TypeName() string       { return "nil" }
func (n *NilNode) Loc() *TokenLocation    { return n.Location }
func (n *NilNode) IsEmpty() bool          { return true }
func (n *NilNode) First() Node            { return n }
func (n *NilNode) Rest() Node             { return &ListNode{} }
func (n *NilNode) Equals(other Node) bool {
	if _, ok := other.(*NilNode); ok {
		return true
	}
	return false
}
func (n *NilNode) Append(other Coll) Coll {
	return other
}
func (n *NilNode) Cons(elem Node) Coll {
	return NewListNode([]Node{elem})
}

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

////////// CharNode

type CharNode struct {
	Value      rune
	annotation Node
	Location   *TokenLocation
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
func (cn *CharNode) Children() []Node     { return nil }
func (cn *CharNode) isExpr() bool         { return true }
func (cn *CharNode) Annotation() Node     { return cn.annotation }
func (cn *CharNode) SetAnnotation(n Node) { cn.annotation = n }
func (cn *CharNode) Equals(n Node) bool   { return cn.Value == asChar(n).Value }
func (cn *CharNode) TypeName() string     { return "char" }
func (cn *CharNode) Loc() *TokenLocation  { return cn.Location }

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

type ListNode struct {
	Nodes      []Node
	annotation Node
	Location   *TokenLocation
}

func NewListNode(nodes []Node) *ListNode {
	return &ListNode{Nodes: nodes}
}

func (l *ListNode) String() string {
	raw := "(" + strings.Join(nodesToStrings(l.Nodes), " ") + ")"
	return displayAnnotation(l, raw)
}

func (l *ListNode) Children() []Node     { return l.Nodes }
func (l *ListNode) isExpr() bool         { return true }
func (l *ListNode) Annotation() Node     { return l.annotation }
func (l *ListNode) SetAnnotation(n Node) { l.annotation = n }
func (l *ListNode) TypeName() string     { return "list" }
func (l *ListNode) Loc() *TokenLocation  { return l.Location }
func (l *ListNode) Equals(n Node) bool {
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
func (l *ListNode) Append(other Coll) Coll {
	if other.IsEmpty() {
		return l
	} else {
		return NewListNode(append(l.Nodes, other.Children()...))
	}
}
func (l *ListNode) Cons(elem Node) Coll {
	return NewListNode(append([]Node{elem}, l.Nodes...))
}
func (l *ListNode) IsEmpty() bool {
	return len(l.Nodes) == 0
}
func (l *ListNode) First() Node {
	if len(l.Nodes) == 0 {
		return &NilNode{}
	}
	return l.Nodes[0]
}
func (l *ListNode) Rest() Node {
	if len(l.Nodes) == 0 {
		return l
	}
	return NewListNode(l.Nodes[1:])
}

////////// Helpers

func asStringNode(n Node) *StringNode {
	if result, ok := n.(*StringNode); ok {
		return result
	}
	return &StringNode{}
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
func asList(n Node) *ListNode {
	if result, ok := n.(*ListNode); ok {
		return result
	}
	return &ListNode{}
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
