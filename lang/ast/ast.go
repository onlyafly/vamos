package ast

import (
	"fmt"
	"strconv"
	"strings"
)

////////// Module

type Module struct {
	Nodes []Node
}

func (m *Module) String() (result string) {
	for i, n := range m.Nodes {
		result += n.String()

		if i < len(m.Nodes)-1 {
			result += "\n\n"
		}
	}
	return
}

////////// Node

type Node interface {
	fmt.Stringer
	Children() []Node
	//TODO Pos() int
}

////////// BasicNode

type BasicNode interface {
	Node
	Annotation() Node
	SetAnnotation(n Node)
}

func displayAnnotation(bn BasicNode, rawRepresentation string) string {
	if bn.Annotation() != nil {
		return "^" + bn.Annotation().String() + " " + rawRepresentation
	}

	return rawRepresentation
}

////////// Expression

type Stmt interface {
	Node
	isStmt() bool
}

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
	raw := "(" + strings.Join(stringNodes(this.Nodes), " ") + ")"
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
		strings.Join(stringNodes(self.Body), " ") +
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

func toNumberValue(exp Expr) float64 {
	switch value := exp.(type) {
	case *Number:
		return value.Value
	}

	panic("Expression is not a number: " + exp.String())
}

func toSymbolValue(exp Expr) string {
	switch value := exp.(type) {
	case *Symbol:
		return value.Name
	}

	panic("Expression is not a symbol: " + exp.String())
}

func stringNodes(nodes []Node) []string {
	return nodesToStrings(nodes, func(n Node) string { return n.String() })
}

func nodesToStrings(nodes []Node, convert func(n Node) string) []string {
	strings := make([]string, len(nodes))
	for i, node := range nodes {
		strings[i] = convert(node)
	}
	return strings
}
