package lang

import (
	"fmt"
	"strconv"
	"strings"
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
	//TODO Pos() int
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

////////// Symbol

type Symbol struct {
	Name       string
	annotation Node
}

func (s *Symbol) String() string       { return displayAnnotation(s, s.Name) }
func (s *Symbol) Children() []Node     { return nil }
func (s *Symbol) isExpr() bool         { return true }
func (s *Symbol) Annotation() Node     { return s.annotation }
func (s *Symbol) SetAnnotation(n Node) { s.annotation = n }
func (s *Symbol) Equals(n Node) bool   { return s.Name == asSymbol(n).Name }

////////// Number

type Number struct {
	Value      float64
	annotation Node
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

////////// List

type List struct {
	Nodes      []Node
	annotation Node
}

func (l *List) String() string {
	raw := "(" + strings.Join(nodesToStrings(l.Nodes), " ") + ")"
	return displayAnnotation(l, raw)
}

func (l *List) Children() []Node     { return l.Nodes }
func (l *List) isExpr() bool         { return true }
func (l *List) Annotation() Node     { return l.annotation }
func (l *List) SetAnnotation(n Node) { l.annotation = n }
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

////////// Primitive

type Primitive struct {
	Name  string
	Value primitiveFunction
}

func (p *Primitive) String() string {
	return "#primitive<" + p.Name + ">"
}

func (p *Primitive) Children() []Node { return nil }
func (p *Primitive) isExpr() bool     { return true }
func (p *Primitive) Equals(n Node) bool {
	panicEvalError("Cannot compare the values of primitive procedures: " +
		p.String() + " and " + n.String())
	return false
}

////////// Function

type Function struct {
	Name       string
	Parameters []Node
	Body       Node
	ParentEnv  Env
}

func (f *Function) String() string {
	return "#function<" + f.Name + ">"
}

func (f *Function) Children() []Node { return nil }
func (f *Function) isExpr() bool     { return true }
func (f *Function) Equals(n Node) bool {
	panicEvalError("Cannot compare the values of functions: " +
		f.String() + " and " + n.String())
	return false
}

////////// Helpers

func asSymbol(n Node) *Symbol {
	if result, ok := n.(*Symbol); ok {
		return result
	}
	panicEvalError("Expected symbol: " + n.String())
	return nil
}
func asNumber(n Node) *Number {
	if result, ok := n.(*Number); ok {
		return result
	}
	panicEvalError("Expected number: " + n.String())
	return nil
}
func asList(n Node) *List {
	if result, ok := n.(*List); ok {
		return result
	}
	panicEvalError("Expected list: " + n.String())
	return nil
}

func toListValue(n Node) *List {
	switch value := n.(type) {
	case *List:
		return value
	}

	panicEvalError("Expression is not a list: " + n.String())
	return nil
}

func toNumberValue(n Node) float64 {
	switch value := n.(type) {
	case *Number:
		return value.Value
	}

	panicEvalError("Expression is not a number: " + n.String())
	return 0.0
}

func toSymbolValue(exp Expr) string {
	switch value := exp.(type) {
	case *Symbol:
		return value.Name
	}

	panicEvalError("Expression is not a symbol: " + exp.String())
	return ""
}

func toBooleanValue(n Node) bool {
	switch value := n.(type) {
	case *Symbol:
		switch value.Name {
		case "true":
			return true
		case "false":
			return false
		}
	}

	panicEvalError("Non-boolean in boolean context: " + n.String())
	return false
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
