package parsing

import (
	"testing"
	"vamos/lang"
	"vamos/lang/ast"
)

func TestParseAtom(t *testing.T) {
	result1, _ := parseAtom("fred").(*ast.Symbol)

	lang.CheckEqualString(t, "fred", result1.Name)

	result2, _ := parseAtom("1").(*ast.Number)
	lang.CheckEqualFloat(t, 1, result2.Value)

	result3, _ := parseAtom("2.4").(*ast.Number)
	lang.CheckEqualFloat(t, 2.4, result3.Value)
}

func TestParse(t *testing.T) {
	result := Parse("(defn init ()  (print 42))")

	lang.CheckEqualString(t, "(defn init () (print 42))", result.String())
}
