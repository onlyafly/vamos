package scanning

import (
	"testing"
	"vamos/lang"
)

func TestScan(t *testing.T) {
	_, tokens := Scan("tester1", "(1 2 3)")

	lang.CheckEqualStringer(t, "(", <-tokens)
	lang.CheckEqualStringer(t, "1", <-tokens)
	lang.CheckEqualStringer(t, "2", <-tokens)
	lang.CheckEqualStringer(t, "3", <-tokens)
	lang.CheckEqualStringer(t, ")", <-tokens)
	lang.CheckEqualStringer(t, "EOF", <-tokens)

	_, tokens = Scan("tester2", "(abc ab2? 3.5)")

	lang.CheckEqualStringer(t, "(", <-tokens)
	lang.CheckEqualStringer(t, "abc", <-tokens)
	lang.CheckEqualStringer(t, "ab2?", <-tokens)
	lang.CheckEqualStringer(t, "3.5", <-tokens)
	lang.CheckEqualStringer(t, ")", <-tokens)
	lang.CheckEqualStringer(t, "EOF", <-tokens)
}
