package testhelp

import (
	"fmt"
	"testing"
)

func CheckEqualStringer(t *testing.T, expected, actual interface{}) {
	e := fmt.Sprintf("%v", expected)
	a := fmt.Sprintf("%v", actual)

	if e != a {
		t.Errorf("Expected <%v>, got <%v>", e, a)
	}
}

func CheckEqualString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Errorf("Expected <%v>, got <%v>", expected, actual)
	}
}

func CheckEqualInt(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected <%v>, got <%v>", expected, actual)
	}
}

func CheckEqualFloat(t *testing.T, expected, actual float64) {
	if expected != actual {
		t.Errorf("Expected <%v>, got <%v>", expected, actual)
	}
}
