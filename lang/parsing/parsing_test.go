package parsing

import (
	"testing"
	"vamos/lang"
	"vamos/lang/scanning"
)

func TestParseAtom(t *testing.T) {
	result1 := parseSymbol(scanning.Token{Value: "fred"})
	lang.CheckEqualString(t, "fred", result1.Name)

	errors := NewParserErrorList()

	result2 := parseNumber(scanning.Token{Value: "1"}, &errors)
	lang.CheckEqualFloat(t, 1, result2.Value)

	result3 := parseNumber(scanning.Token{Value: "2.4"}, &errors)
	lang.CheckEqualFloat(t, 2.4, result3.Value)
}

func TestParse(t *testing.T) {
	result, _ := Parse("(defn init ()  (print 42))")

	lang.CheckEqualString(t, "(defn init () (print 42))", result.String())
}

func TestParse_SymbolAnnotatingSymbol(t *testing.T) {
	result, _ := Parse("(defn ^sample init ()  (print 42))")

	lang.CheckEqualString(t, "(defn ^sample init () (print 42))", result.String())
}
