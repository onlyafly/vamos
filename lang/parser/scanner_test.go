package parser

import (
	"fmt"
	"github.com/onlyafly/vamos/testhelp"
	"testing"
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

	_, tokens = Scan("tester3", "\\a")
	testhelp.CheckEqualStringer(t, "\\a", <-tokens)

	fmt.Printf("START\n")
	_, tokens = Scan("tester4", "(list \\a)")
	testhelp.CheckEqualStringer(t, "(", <-tokens)
	testhelp.CheckEqualStringer(t, "list", <-tokens)
	testhelp.CheckEqualStringer(t, "\\a", <-tokens)
	testhelp.CheckEqualStringer(t, ")", <-tokens)
	fmt.Printf("END\n")
}
