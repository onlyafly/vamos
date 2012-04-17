package lang

import ()

////////// Evaluation

func Eval(e Env, n Node) Node {
	switch value := n.(type) {
	case *Number:
		return value
	case *Symbol:
		return e.Get(value.Name)
	case *List:
		return evalList(e, value)
	default:
		panic("Unknown form to evaluate: " + value.String())
	}

	return &Symbol{Name: "nil"}
}

func evalList(e Env, l *List) Node {
	elements := l.Nodes
	proc := elements[0]
	args := elements[1:]

	switch value := proc.(type) {
	case *Symbol:
		switch value.Name {
		case "+":
			result := toNumberValue(Eval(e, args[0])) + toNumberValue(Eval(e, args[1]))
			return &Number{Value: result}
		case "def":
			name := toSymbolName(args[0])
			e.Set(name, Eval(e, args[1]))
			return &Symbol{Name: "nil"}
		case "quote":
			return args[0]
		default:
			panic("Unknown function to evaluate: " + value.Name)
		}
	default:
		panic("First item in list not a symbol: " + value.String())
	}

	return &Symbol{Name: "nil"}
}

func toSymbolName(n Node) string {
	switch value := n.(type) {
	case *Symbol:
		return value.Name
	}

	panic("Not a symbol: " + n.String())
}

/*
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
*/
