package ast

import (
	"fmt"
	"strconv"
	"strings"
)

////////// Module

type Module struct {
	Expressions []Expression
}

func (m *Module) String() (result string) {
	for i, e := range m.Expressions {
		result += e.String()

		if i < len(m.Expressions)-1 {
			result += "\n\n"
		}
	}
	return
}

////////// Node

type Node interface {
	Pos() int
}

////////// Expression

type Expression interface {
	String() string
}

////////// Symbol

type Symbol struct {
	Name string
}

func NewSymbol(name string) *Symbol {
	return &Symbol{name}
}

func (self *Symbol) String() string {
	return self.Name
}

////////// Number

type Number struct {
	Value float64
}

func NewNumber(value float64) *Number {
	return &Number{value}
}

func (self *Number) String() string {
	return strconv.FormatFloat(
		self.Value,
		'f',
		-1,
		64)
}

////////// List

type List struct {
	Value []Expression
}

func NewList(value []Expression) *List {
	return &List{value}
}

func (self *List) String() string {
	return "(" + strings.Join(stringExpressions(self.Value), " ") + ")"
}

////////// Function Definition

type FunctionDefinition struct {
	Name      *Symbol
	Arguments *List
	Body      []Expression
}

func NewFunctionDefinition(name *Symbol, args *List, body []Expression) *FunctionDefinition {
	return &FunctionDefinition{name, args, body}
}

func (self *FunctionDefinition) String() string {
	return "(defn " +
		self.Name.String() + " " + self.Arguments.String() + " " +
		strings.Join(stringExpressions(self.Body), " ") +
		")"
}

////////// Package Definition

type PackageDefinition struct {
	Name *Symbol
}

func NewPackageDefinition(name *Symbol) *PackageDefinition {
	return &PackageDefinition{name}
}

func (self *PackageDefinition) String() string {
	return fmt.Sprintf("(package %v)", self.Name)
}

////////// Helpers

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

func expressionsToStrings(es []Expression, convert func(e Expression) string) []string {
	strings := make([]string, len(es))
	for i, exp := range es {
		strings[i] = convert(exp)
	}
	return strings
}
