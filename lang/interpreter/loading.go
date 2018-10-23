package interpreter

import (
	"fmt"
	"os"

	"github.com/onlyafly/vamos/lang/ast"
	"github.com/onlyafly/vamos/lang/parser"
)

// ParseEvalPrint parses, evals, and prints
func ParseEvalPrint(env Env, input string, readLine func() string, sourceName string, printResult bool) {
	if result, err := ParseEval(env, input, readLine, sourceName); err == nil {
		// Can be null if nothing was entered
		if result != nil && printResult {
			fmt.Println(result.String())
		}
	} else {
		fmt.Println(err.Error())
	}
}

// ParseEval parses and evals
func ParseEval(env Env, input string, readLine func() string, sourceName string) (ast.Node, error) {
	defer func() {
		// Some non-application triggered panic has occurred
		if e := recover(); e != nil {
			fmt.Printf("Host environment error: %v\n", e)
			panic(e)
		}
	}()

	nodes, parseErrors := parser.Parse(input, sourceName)

	if parseErrors != nil {
		fmt.Println(parseErrors.String())
	}

	var result ast.Node
	var evalError error
	for _, n := range nodes {
		result, evalError = Eval(env, n, os.Stdout, readLine)
		if evalError != nil {
			break
		}
	}

	if evalError == nil {
		return result, nil
	}

	return nil, evalError
}
