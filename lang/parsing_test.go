package lang

import (
	"testing"
)

func TestTokenize(t *testing.T) {
	result := tokenize("(a b 10)")
	checkEqualString(t, "(", result.pop())
	checkEqualString(t, "a", result.pop())
	checkEqualString(t, "b", result.pop())
	checkEqualString(t, "10", result.pop())
	checkEqualString(t, ")", result.pop())
}

func TestParseAtom(t *testing.T) {
	result1, _ := parseAtom("fred").(*Symbol)

	checkEqualString(t, "fred", result1.Name)

	result2, _ := parseAtom("1").(*Number)
	checkEqualFloat(t, 1, result2.Value)

	result3, _ := parseAtom("2.4").(*Number)
	checkEqualFloat(t, 2.4, result3.Value)
}

func TestParse(t *testing.T) {
	result := Parse("(defn init () (print 42))")

	checkEqualString(t, "<DEFN init ARGS=() BODY=(print 42)>", result.Probe())
}
