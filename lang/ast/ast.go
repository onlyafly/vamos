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
	//TODO Pos() int
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
	Name string
}

func (self *Symbol) String() string {
	return self.Name
}

func (self *Symbol) isExpr() bool { return true }

////////// Number

type Number struct {
	Value float64
}

func (self *Number) String() string {
	return strconv.FormatFloat(
		self.Value,
		'f',
		-1,
		64)
}

func (self *Number) isExpr() bool { return true }

////////// List

type List struct {
	Nodes []Node
}

func (self *List) String() string {
	return "(" + strings.Join(stringNodes(self.Nodes), " ") + ")"
}

func (self *List) isExpr() bool { return true }

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

func (self *FunctionDecl) isDecl() bool { return true }

////////// Package Declaration

type PackageDecl struct {
	Name *Symbol
}

func (self *PackageDecl) String() string {
	return fmt.Sprintf("(package %v)", self.Name)
}

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
