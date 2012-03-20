package compiling

import (
	"fmt"
	"vamos/lang/ast"
)

////////// Evaluation

func Compile(m *ast.Module) (result string) {
	for i, e := range m.Expressions {
		switch value := e.(type) {
		case *ast.FunctionDefinition:
			result += CompileFunctionDefinition(value)
		case *ast.PackageDefinition:
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

func CompileFunctionDefinition(fd *ast.FunctionDefinition) (r string) {
	r = "func " + fd.Name.String() + "() {\n"

	for _, e := range fd.Body {
		r += "\t" + CompileBodyExpression(e) + "\n"
	}

	r += "}"
	return
}

func CompilePackageDefinition(pd *ast.PackageDefinition) (r string) {
	r = fmt.Sprintf("package %v", pd.Name)
	return
}

func CompileBodyExpression(e ast.Expression) string {
	switch value := e.(type) {
	case *ast.List:
		return CompileFunctionCall(value)
	default:
		panic("Cannot compile: " + value.String())
	}

	return ""
}

func CompileFunctionCall(list *ast.List) (r string) {
	r = fmt.Sprintf("%v(", list.Value[0])
	for i, e := range list.Value[1:] {
		r += e.String()
		if i < len(list.Value[1:])-1 {
			r += ", "
		}
	}
	r += ")"
	return
}
