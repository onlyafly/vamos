package lang

import (
	"testing"

	"../util"
)

func TestScan(t *testing.T) {
	_, tokens := Scan("tester1", "(1 2 3)")

	util.CheckEqualStringer(t, "(", <-tokens)
	util.CheckEqualStringer(t, "1", <-tokens)
	util.CheckEqualStringer(t, "2", <-tokens)
	util.CheckEqualStringer(t, "3", <-tokens)
	util.CheckEqualStringer(t, ")", <-tokens)
	util.CheckEqualStringer(t, "EOF", <-tokens)

	_, tokens = Scan("tester2", "(abc ab2? 3.5)")

	util.CheckEqualStringer(t, "(", <-tokens)
	util.CheckEqualStringer(t, "abc", <-tokens)
	util.CheckEqualStringer(t, "ab2?", <-tokens)
	util.CheckEqualStringer(t, "3.5", <-tokens)
	util.CheckEqualStringer(t, ")", <-tokens)
	util.CheckEqualStringer(t, "EOF", <-tokens)
}
