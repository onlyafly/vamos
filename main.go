package main

import (
	"flag"
	"fmt"
	"vamos/lang"
	"vamos/util"
)

const (
	version     = `0.1.1`
	versionDate = `2015-01-29`
)

func main() {
	fmt.Printf("Vamos %s (%s)\n", version, versionDate)

	env := lang.NewTopLevelMapEnv()

	// Loading of files

	fileName := flag.String("l", "", "load a file at startup")
	flag.Parse()

	if fileName != nil && len(*fileName) > 0 {
		content, _ := util.ReadFile(*fileName)
		parseEval(env, content)
	}

	// REPL

	for {
		fmt.Print("> ")
		input := util.ReadLine()
		parseEval(env, input)
	}
}

func parseEval(env lang.Env, input string) {
	nodes, parseErrors := lang.Parse(input)

	if parseErrors != nil {
		fmt.Println(parseErrors.String())
	}

	var result lang.Node
	var evalError error
	for _, n := range nodes {
		result, evalError = lang.Eval(env, n)
		if evalError != nil {
			break
		}
	}

	var actual string
	if evalError == nil {
		actual = result.String()
	} else {
		actual = evalError.Error()
	}

	fmt.Println(actual)
}
