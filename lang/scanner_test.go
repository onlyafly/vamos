package lang

import (
	"testing"
)

func TestScan(t *testing.T) {
	_, tokens := Scan("tester1", "(1 2 3)")

	checkEqualStringer(t, "(", <-tokens)
	checkEqualStringer(t, "1", <-tokens)
	checkEqualStringer(t, "2", <-tokens)
	checkEqualStringer(t, "3", <-tokens)
	checkEqualStringer(t, ")", <-tokens)
	checkEqualStringer(t, "EOF", <-tokens)

	_, tokens = Scan("tester2", "(abc ab2? 3.5)")

	checkEqualStringer(t, "(", <-tokens)
	checkEqualStringer(t, "abc", <-tokens)
	checkEqualStringer(t, "ab2?", <-tokens)
	checkEqualStringer(t, "3.5", <-tokens)
	checkEqualStringer(t, ")", <-tokens)
	checkEqualStringer(t, "EOF", <-tokens)
}
