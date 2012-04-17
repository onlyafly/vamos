package lang

import (
	"fmt"
	"strconv"
	"strings"
)

////////// Slice of Nodes

type Nodes []Node

func (ns Nodes) String() string {
	return strings.Join(nodesToStrings([]Node(ns)), "\n")
}

////////// Node

type Node interface {
	fmt.Stringer
	Children() []Node
	//TODO Pos() int
}

////////// AnnotatedNode

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

type Expr interface {
	Node
	isExpr() bool
}

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

////////// Number

type Number struct {
	Value      float64
	annotation Node
}

func (n *Number) String() string {
	rep := strconv.FormatFloat(
		n.Value,
		'f',
		-1,
		64)

	return displayAnnotation(n, rep)
}

func (this *Number) Children() []Node     { return nil }
func (this *Number) isExpr() bool         { return true }
func (this *Number) Annotation() Node     { return this.annotation }
func (this *Number) SetAnnotation(n Node) { this.annotation = n }

////////// List

type List struct {
	Nodes      []Node
	annotation Node
}

func (this *List) String() string {
	raw := "(" + strings.Join(nodesToStrings(this.Nodes), " ") + ")"
	return displayAnnotation(this, raw)
}

func (self *List) Children() []Node     { return self.Nodes }
func (self *List) isExpr() bool         { return true }
func (this *List) Annotation() Node     { return this.annotation }
func (this *List) SetAnnotation(n Node) { this.annotation = n }

////////// Function Decl

type FunctionDecl struct {
	Name      *Symbol
	Arguments *List
	Body      []Node
}

func (self *FunctionDecl) String() string {
	return "(defn " +
		self.Name.String() + " " + self.Arguments.String() + " " +
		strings.Join(nodesToStrings(self.Body), " ") +
		")"
}

func (self *FunctionDecl) Children() []Node { return nil }

func (self *FunctionDecl) isDecl() bool { return true }

////////// Package Declaration

type PackageDecl struct {
	Name *Symbol
}

func (self *PackageDecl) String() string {
	return fmt.Sprintf("(package %v)", self.Name)
}

func (self *PackageDecl) Children() []Node { return []Node{self.Name} }

func (self *PackageDecl) isDecl() bool { return true }

////////// Helpers

func toNumberValue(n Node) float64 {
	switch value := n.(type) {
	case *Number:
		return value.Value
	}

	panic("Expression is not a number: " + n.String())
}

func toSymbolValue(exp Expr) string {
	switch value := exp.(type) {
	case *Symbol:
		return value.Name
	}

	panic("Expression is not a symbol: " + exp.String())
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
