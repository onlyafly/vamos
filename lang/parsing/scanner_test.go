package parsing

import (
	"testing"
	"vamos/testhelp"
)

func TestScan(t *testing.T) {
	_, tokens := Scan("tester1", "(1 2 3)")

	testhelp.CheckEqualStringer(t, "(", <-tokens)
	testhelp.CheckEqualStringer(t, "1", <-tokens)
	testhelp.CheckEqualStringer(t, "2", <-tokens)
	testhelp.CheckEqualStringer(t, "3", <-tokens)
	testhelp.CheckEqualStringer(t, ")", <-tokens)
	testhelp.CheckEqualStringer(t, "EOF", <-tokens)

	_, tokens = Scan("tester2", "(abc ab2? 3.5)")

	testhelp.CheckEqualStringer(t, "(", <-tokens)
	testhelp.CheckEqualStringer(t, "abc", <-tokens)
	testhelp.CheckEqualStringer(t, "ab2?", <-tokens)
	testhelp.CheckEqualStringer(t, "3.5", <-tokens)
	testhelp.CheckEqualStringer(t, ")", <-tokens)
	testhelp.CheckEqualStringer(t, "EOF", <-tokens)
}
