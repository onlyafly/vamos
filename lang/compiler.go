package lang

import (
	"fmt"	
)

////////// Evaluation

func Compile(m *Module) (result string) {
	for i, e := range m.Expressions {
		switch value := e.(type) {
		case *FunctionDefinition:
			result += CompileFunctionDefinition(value)
		case *PackageDefinition:
			result += CompilePackageDefinition(value)
		default:
			panic("Cannot compile: " + value.String())
		}

		if i < len(m.Expressions)-1 {
			result += "\n\n"
		}
	}

	return
}

func CompileFunctionDefinition(fd *FunctionDefinition) (r string) {
	r = "func " + fd.Name.String() + "() {\n"

	for _, e := range fd.Body {
		r += "\t" + CompileBodyExpression(e) + "\n"
	}

	r += "}"
	return
}

func CompilePackageDefinition(pd *PackageDefinition) (r string) {
	r = fmt.Sprintf("package %v", pd.Name)
	return
}

func CompileBodyExpression(e Expression) string {
	switch value := e.(type) {
	case *List:
		return CompileFunctionCall(value)
	default:
		panic("Cannot compile: " + value.String())
	}

	return ""
}

func CompileFunctionCall(list *List) (r string) {
	r = fmt.Sprintf("%v(", list.Value[0])
	for i, e := range list.Value[1:] {
		r += e.String()
		if i < len(list.Value[1:]) - 1 {
			r += ", "
		}
	}
	r += ")"
	return
}
