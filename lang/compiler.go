package lang

import ()

////////// Evaluation

func Compile(e Expression) string {
	
	result := "package main\n\n"

	switch value := e.(type) {
	case *FunctionDefinition:
		result += CompileFunctionDefinition(value)
	default:
		panic("Cannot compile: " + value.Probe())
	}

	return result
}

func CompileFunctionDefinition(fd *FunctionDefinition) (r string) {
	r = "func " + fd.Name.String() + "() {\n"
	
	for _, e := range fd.Body {
		r += "\t" + CompileBodyExpression(e) + "\n"
	}

	r += "}"
	return
}

func CompileBodyExpression(e Expression) string {
	switch value := e.(type) {
	case *List:
		return CompileFunctionCall(value)
	default:
		panic("Cannot compile: " + value.Probe())
	}

	return ""
}

func CompileFunctionCall(list *List) (r string) {
	r = list.Value[0].String() + "(" + list.Value[1].String() + ")"
	return
}
