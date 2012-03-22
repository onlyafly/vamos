package compiling

import (
	"fmt"
	"vamos/lang/ast"
)

////////// Evaluation

func Compile(m *ast.Module) (result string) {
	for i, n := range m.Nodes {
		switch value := n.(type) {
		case *ast.FunctionDecl:
			result += CompileFunctionDecl(value)
		case *ast.PackageDecl:
			result += CompilePackageDecl(value)
		default:
			panic("Cannot compile: " + value.String())
		}

		if i < len(m.Nodes)-1 {
			result += "\n\n"
		}
	}

	return
}

func CompileFunctionDecl(fd *ast.FunctionDecl) (r string) {
	r = "func " + fd.Name.String() + "() {\n"

	for _, n := range fd.Body {
		r += "\t" + CompileBodyNode(n) + "\n"
	}

	r += "}"
	return
}

func CompilePackageDecl(pd *ast.PackageDecl) (r string) {
	r = fmt.Sprintf("package %v", pd.Name)
	return
}

func CompileBodyNode(n ast.Node) string {
	switch value := n.(type) {
	case *ast.List:
		return CompileFunctionCall(value)
	default:
		panic("Cannot compile: " + value.String())
	}

	return ""
}

func CompileFunctionCall(list *ast.List) (r string) {
	r = fmt.Sprintf("%v(", list.Nodes[0])
	for i, e := range list.Nodes[1:] {
		r += e.String()
		if i < len(list.Nodes[1:])-1 {
			r += ", "
		}
	}
	r += ")"
	return
}
