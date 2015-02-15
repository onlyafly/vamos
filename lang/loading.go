package lang

import (
	"fmt"
	"os"
)

func ParseEvalPrint(env Env, input string, sourceName string, printResult bool) {
	if result, err := ParseEval(env, input, sourceName); err == nil {
		// Can be null if nothing was entered
		if result != nil && printResult {
			fmt.Println(result.String())
		}
	} else {
		fmt.Println(err.Error())
	}
}

func ParseEval(env Env, input string, sourceName string) (Node, error) {
	defer func() {
		// Some non-application triggered panic has occurred
		if e := recover(); e != nil {
			fmt.Printf("Host environment error: %v\n", e)
			panic(e)
		}
	}()

	nodes, parseErrors := Parse(input, sourceName)

	if parseErrors != nil {
		fmt.Println(parseErrors.String())
	}

	var result Node
	var evalError error
	for _, n := range nodes {
		result, evalError = Eval(env, n, os.Stdout)
		if evalError != nil {
			break
		}
	}

	if evalError == nil {
		return result, nil
	} else {
		return nil, evalError
	}
}
