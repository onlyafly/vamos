package lang

import (
	"strconv"
	"strings"
)

type Expression interface {
	Probe() string
	String() string
}

type Symbol struct {
	Name string
}

func NewSymbol(name string) *Symbol {
	return &Symbol{name}
}

func (self *Symbol) Probe() string {
	return self.Name
}

func (self *Symbol) String() string {
	return self.Name
}

type Number struct {
	Value float64
}

func NewNumber(value float64) *Number {
	return &Number{value}
}

func (self *Number) Probe() string {
	return self.String()
}

func (self *Number) String() string {
	return strconv.FormatFloat(
		self.Value,
		'f',
		-1,
		64)
}

type List struct {
	Value []Expression
}

func NewList(value []Expression) *List {
	return &List{value}
}

func (self *List) Probe() string {
	return "(" + strings.Join(probeExpressions(self.Value), " ") + ")"
	return "list"
}

func (self *List) String() string {
	return "(" + strings.Join(stringExpressions(self.Value), " ") + ")"
}

type FunctionDefinition struct {
	Name      *Symbol
	Arguments *List
	Body      []Expression
}

func NewFunctionDefinition(name *Symbol, args *List, body []Expression) *FunctionDefinition {
	return &FunctionDefinition{name, args, body}
}

func (self *FunctionDefinition) Probe() string {
	return "<DEFN " + self.Name.Probe() +
		" ARGS=" + self.Arguments.Probe() +
		" BODY=" +
		strings.Join(probeExpressions(self.Body), " ") +
		">"
}

func (self *FunctionDefinition) String() string {
	return "(defn " +
		self.Name.String() + " " + self.Arguments.String() + " " +
		strings.Join(stringExpressions(self.Body), " ") +
		")"
}

func toNumberValue(exp Expression) float64 {
	switch value := exp.(type) {
	case *Number:
		return value.Value
	}

	panic("Expression is not a number: " + exp.String())
}

func toSymbolValue(exp Expression) string {
	switch value := exp.(type) {
	case *Symbol:
		return value.Name
	}

	panic("Expression is not a symbol: " + exp.String())
}

func stringExpressions(es []Expression) []string {
	return expressionsToStrings(es, stringExpression)
}

func stringExpression(e Expression) string {
	return e.String()
}

func probeExpressions(es []Expression) []string {
	return expressionsToStrings(es, probeExpression)
}

func probeExpression(e Expression) string {
	return e.Probe()
}

func expressionsToStrings(es []Expression, convert func(e Expression)string) []string {
	strings := make([]string, len(es))
	for i, exp := range es {
		strings[i] = convert(exp)
	}
	return strings
}
