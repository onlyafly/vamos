package lang

import ()

////////// Evaluation

func Eval(e Env, n Node) (result Node, err error) {
	defer func() {
		if e := recover(); e != nil {
			result = nil
			switch errorValue := e.(type) {
			case EvalError:
				err = errorValue
				return
			default:
				panic(errorValue)
			}
		}
	}()

	return evalNode(e, n), nil
}

func evalEachNode(e Env, ns []Node) []Node {
	result := make([]Node, len(ns))
	for i, n := range ns {
		result[i] = evalNode(e, n)
	}
	return result
}

func evalNode(e Env, n Node) Node {
	switch value := n.(type) {
	case *Number:
		return value
	case *Symbol:
		return e.Get(value.Name)
	case *List:
		return evalList(e, value)
	default:
		panicEvalError("Unknown form to evaluate: " + value.String())
	}

	return &Symbol{Name: "nil"}
}

func evalList(e Env, l *List) Node {
	elements := l.Nodes
	head := elements[0]
	args := elements[1:]

	switch value := head.(type) {
	case *Symbol:
		switch value.Name {
		case "def":
			name := toSymbolName(args[0])
			e.Set(name, evalNode(e, args[1]))
			return &Symbol{Name: "nil"}
		case "if":
			predicate := toBooleanValue(evalNode(e, args[0]))
			if predicate {
				return evalNode(e, args[1])
			} else {
				return evalNode(e, args[2])
			}
		case "quote":
			return args[0]
		}
	}

	headValue := evalNode(e, head)

	switch value := headValue.(type) {
	case *Primitive:
		f := value.Value
		return f(evalEachNode(e, args))
	default:
		panicEvalError("First item in list not a function: " + value.String())
	}

	return &Symbol{Name: "nil"}
}

func panicEvalError(s string) {
	panic(EvalError(s))
}

func toSymbolName(n Node) string {
	switch value := n.(type) {
	case *Symbol:
		return value.Name
	}

	panic("Not a symbol: " + n.String())
}
