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
	Equals(Node) bool
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

func (this *Symbol) String() string       { return displayAnnotation(this, this.Name) }
func (this *Symbol) Children() []Node     { return nil }
func (this *Symbol) isExpr() bool         { return true }
func (this *Symbol) Annotation() Node     { return this.annotation }
func (this *Symbol) SetAnnotation(n Node) { this.annotation = n }
func (this *Symbol) Equals(n Node) bool   { return this.Name == asSymbol(n).Name }

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
func (this *Number) Equals(n Node) bool   { return this.Value == asNumber(n).Value }

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
func (this *List) Equals(n Node) bool {
	other := asList(n)

	// Compare lengths
	if len(this.Nodes) != len(other.Nodes) {
		return false
	}

	// Compare contents
	for i, v := range this.Nodes {
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

func (this *Primitive) String() string {
	return "#primitive<" + this.Name + ">"
}

func (this *Primitive) Children() []Node { return nil }
func (this *Primitive) isExpr() bool     { return true }
func (this *Primitive) Equals(n Node) bool {
	panicEvalError("Cannot compare the values of primitive procedures: " +
		this.String() + " and " + n.String())
	return false
}

////////// Function

type Function struct {
	Name       string
	Parameters []Node
	Body       Node
	ParentEnv  Env
}

func (this *Function) String() string {
	return "#function<" + this.Name + ">"
}

func (this *Function) Children() []Node { return nil }
func (this *Function) isExpr() bool     { return true }
func (this *Function) Equals(n Node) bool {
	panicEvalError("Cannot compare the values of functions: " +
		this.String() + " and " + n.String())
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
