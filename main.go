package main

import (
	//"flag"
	"fmt"
	"vamos/lang/scanning"
)

func main() {
	fmt.Println("Vamos!")

	_, c := scanning.Scan("x", "(^int 2 3)")
	fmt.Println("result: %v", <-c)
	fmt.Println("result: %v", <-c)
	fmt.Println("result: %v", <-c)
	fmt.Println("result: %v", <-c)
	fmt.Println("result: %v", <-c)

	/*
		//fileName := flag.String("c", "", "compile a file")
		flag.Parse()
		fileName := flag.Arg(0)

		content, _ := lang.ReadFile(fileName)

		ast := lang.Parse(content)
		result := lang.Compile(ast)

		_ = lang.WriteFile("output.go", result)
		fmt.Println("Wrote output.go")
	*/
}
