package parser

import (
	"testing"
	"vamos/testhelp"
)

func TestParseAtom(t *testing.T) {
	errors := NewParserErrorList()

	result1 := parseSymbol(Token{Value: "fred"}, &errors)
	testhelp.CheckEqualString(t, "fred", result1.String())

	result2 := parseNumber(Token{Value: "1"}, &errors)
	testhelp.CheckEqualFloat(t, 1, result2.Value)

	result3 := parseNumber(Token{Value: "2.4"}, &errors)
	testhelp.CheckEqualFloat(t, 2.4, result3.Value)
}

func TestParse(t *testing.T) {
	result, _ := Parse("(defn init ()  (print 42))", "test")

	testhelp.CheckEqualString(t, "((defn init () (print 42)))", result.String())
}

func TestParse_SymbolAnnotatingSymbol(t *testing.T) {
	result, _ := Parse("(defn ^sample init ()  (print 42))", "test")

	testhelp.CheckEqualString(t, "((defn ^sample init () (print 42)))", result.String())
}
