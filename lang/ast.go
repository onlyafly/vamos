package lang

import (
	"strings"
	"strconv"
)

type Form interface {
	Type() string
	String() string
}

type Symbol struct {
	Name string
}

func NewSymbol(name string) *Symbol {
	return &Symbol{name}
}

func (self *Symbol) Type() string {
	return "symbol"
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

func (self *Number) Type() string {
	return "number"
}

func (self *Number) String() string {
	return strconv.FormatFloat(
		self.Value,
		'f',
		-1,
		64)
}

type List struct {
	Value []Form
}

func NewList(value []Form) *List {
	return &List{value}
}

func (self *List) Type() string {
	return "list"
}

func (self *List) String() string {
	forms := make([]string, len(self.Value))
	for i, form := range self.Value {
		forms[i] = form.String()
	}
	return "(" + strings.Join(forms, " ") + ")"
}

func toNumberValue(form Form) float64 {
	switch value := form.(type) {
	case *Number:
		return value.Value
	}
	
	panic("Form is not a number: " + form.String())
}

func toSymbolValue(form Form) string {
	switch value := form.(type) {
	case *Symbol:
		return value.Name
	}
	
	panic("Form is not a symbol: " + form.String())
}
