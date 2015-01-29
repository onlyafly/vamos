package main

import (
	"flag"
	"fmt"
	"vamos/lang"
	"vamos/util"
)

func main() {
	fmt.Println("Vamos 0.1.1 (2015-01-29)")

	/*
			fmt.Fprint

			_, c := scanning.Scan("x", "(^int 2 3)")
			fmt.Println("result: %v", <-c)
			fmt.Println("result: %v", <-c)
			fmt.Println("result: %v", <-c)
			fmt.Println("result: %v", <-c)
			fmt.Println("result: %v", <-c)

				//fileName := flag.String("c", "", "compile a file")
				flag.Parse()
				fileName := flag.Arg(0)

		 content, _ := lang.ReadFile(fileName)

				ast := lang.Parse(content)
				result := lang.Compile(ast)

				_ = lang.WriteFile("output.go", result)
				fmt.Println("Wrote output.go")
	*/

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
	} else {
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
}
