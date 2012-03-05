package main

import (
	"fmt"
	"vamos/lang"
)

func main() {
	fmt.Println("Vamos!")
	fmt.Println("Press CTRL+C to quit.")

	for {
		fmt.Print("> ")

		input := lang.ReadLine()

		ast := lang.Parse(input)

		env := lang.NewEnv()
		result := lang.Eval(ast, env)
		fmt.Println(result)
	}
}
