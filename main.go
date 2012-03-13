package main

import (
	"fmt"
	"flag"
	"vamos/lang"
)

func main() {
	fmt.Println("Vamos!")

	//fileName := flag.String("c", "", "compile a file")
	flag.Parse()
	fileName := flag.Arg(0)

	content, _ := lang.ReadFile(fileName)

	ast := lang.Parse(content)
	result := lang.Compile(ast)

	_ = lang.WriteFile("output.go", result)
	fmt.Println("Wrote output.go")
}