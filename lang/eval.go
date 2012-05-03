package lang

import ()

////////// Trampoline Support

type packet struct {
	Thunk thunk
	Node  Node
}

func bounce(t thunk) packet {
	return packet{Thunk: t}
}

func respond(n Node) packet {
	return packet{Node: n}
}

type thunk func() packet

func trampoline(currentThunk thunk) Node {
	for currentThunk != nil {
		next := currentThunk()

		if next.Thunk != nil {
			currentThunk = next.Thunk
		} else {
			return next.Node
		}
	}

	return nil
}

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

	startThunk := func() packet {
		return evalNode(false, e, n)
	}

	return trampoline(startThunk), nil
}

func evalEachNode(e Env, ns []Node) []Node {
	result := make([]Node, len(ns))
	for i, n := range ns {
		evalNodeThunk := func() packet {
			return evalNode(false, e, n)
		}
		result[i] = trampoline(evalNodeThunk)
	}
	return result
}

func evalNode(isTail bool, e Env, n Node) packet {

	switch value := n.(type) {
	case *Number:
		return respond(value)
	case *Symbol:
		return respond(e.Get(value.Name))
	case *List:
		return bounce(func() packet { return evalList(isTail, e, value) })
	default:
		panicEvalError("Unknown form to evaluate: " + value.String())
	}

	return respond(&Symbol{Name: "nil"})
}

func evalList(isTail bool, e Env, l *List) packet {
	elements := l.Nodes

	if len(elements) == 0 {
		panicEvalError("Empty list cannot be evaluated: " + l.String())
		return respond(nil)
	}

	head := elements[0]
	args := elements[1:]

	switch value := head.(type) {
	case *Symbol:
		switch value.Name {
		case "def":
			name := toSymbolName(args[0])
			e.Set(name, trampoline(func() packet {
				return evalNode(false, e, args[1])
			}))
			return respond(&Symbol{Name: "nil"})
		case "if":
			predicate := toBooleanValue(trampoline(func() packet {
				return evalNode(false, e, args[0])
			}))
			if predicate {
				return bounce(func() packet {
					return evalNode(isTail, e, args[1])
				})
			} else {
				return bounce(func() packet {
					return evalNode(isTail, e, args[2])
				})
			}
		case "fn":
			return bounce(func() packet {
				return evalFunctionDefinition(e, args)
			})
		case "quote":
			return respond(args[0])
		}
	}

	var headNode Node = trampoline(func() packet {
		return evalNode(false, e, head)
	})

	switch value := headNode.(type) {
	case *Primitive:
		f := value.Value
		return respond(f(evalEachNode(e, args)))
	case *Function:
		arguments := evalEachNode(e, args)

		return bounce(func() packet {
			return evalFunctionApplication(isTail, value, arguments)
		})
	default:
		panicEvalError("First item in list not a function: " + value.String())
	}

	return respond(&Symbol{Name: "nil"})
}

func evalFunctionApplication(isTail bool, f *Function, args []Node) packet {

	var e Env
	e = NewMapEnv(f.Name, f.ParentEnv)

	// TODO
	/*
		print(
			"evalFunctionApplication:\n   name=",
			e.String(), "\n   body=",
			f.Body.String(), "\n   parent=",
			f.ParentEnv.String(), "\n   args=",
			fmt.Sprintf("%v", args), "\n   isTail=",
			fmt.Sprintf("%v", isTail), "\n")
	*/

	// Save arguments into parameters
	for i, arg := range args {
		paramName := toSymbolName(f.Parameters[i])
		e.Set(paramName, arg)
	}

	return bounce(func() packet {
		return evalNode(true, e, f.Body)
	})
}

func evalFunctionDefinition(e Env, args []Node) packet {
	parameterList := args[0]
	parameterNodes := parameterList.Children()

	return respond(&Function{
		Name:       "anonymous",
		Parameters: parameterNodes,
		Body:       args[1],
		ParentEnv:  e,
	})
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
