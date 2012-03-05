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

func TestAtom(t *testing.T) {
	result1, _ := atom("fred").(*Symbol)
	
	checkEqualString(t, "fred", result1.Name)

	result2, _ := atom("1").(*Number)
	checkEqualFloat(t, 1, result2.Value)

	result3, _ := atom("2.4").(*Number)
	checkEqualFloat(t, 2.4, result3.Value)
}

func TestEvalForms(t *testing.T) {
	forms := []Form {
		NewNumber(3.3),
		NewSymbol("test"),
	}
	env := NewEnv()
	result := evalForms(forms, env)
	checkEqualString(t, "3.3", result[0].String())
	checkEqualString(t, "32.4", result[1].String())
}

func TestEval(t *testing.T) {
	env := NewEnv()
	result := eval(parse("(+ 2 3)"), env)
	checkEqualString(t, "5", result.String())
}

func checkEqualString(t *testing.T, expected, actual string) {
	if (expected != actual) {
		t.Errorf("Expected <%v>, got <%v>", expected, actual)
	}
}

func checkEqualInt(t *testing.T, expected, actual int) {
	if (expected != actual) {
		t.Errorf("Expected <%v>, got <%v>", expected, actual)
	}
}

func checkEqualFloat(t *testing.T, expected, actual float64) {
	if (expected != actual) {
		t.Errorf("Expected <%v>, got <%v>", expected, actual)
	}
}
