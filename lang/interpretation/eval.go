package interpretation

import (
	"strconv"
	. "vamos/lang/ast"
)

////////// Trampoline Support

// Packet contains a thunk or a Node.
// A packet is the result of the evaluation of a thunk.
type packet struct {
	Thunk thunk
	Node  Node
}

// Bounce continues the trampolining session by placing a new thunk in the chain.
func bounce(t thunk) packet {
	return packet{Thunk: t}
}

// Respond exits a trampolining session by placing a Node on the end of the
// chain.
func respond(n Node) packet {
	return packet{Node: n}
}

type thunk func() packet

// Trampoline iteratively calls a chain of thunks until there is no next thunk,
// at which point it pulls the resulting Node out of the packet and returns it.
func trampoline(currentThunk thunk) Node {
	for currentThunk != nil {
		nextPacket := currentThunk()

		if nextPacket.Thunk != nil {
			currentThunk = nextPacket.Thunk
		} else {
			return nextPacket.Node
		}
	}

	return nil
}

////////// Evaluation

// Eval evaluates a node in an environment.
func Eval(e Env, n Node) (result Node, err error) {
	defer func() {
		if e := recover(); e != nil {
			result = nil
			switch errorValue := e.(type) {
			case *EvalError:
				err = errorValue
				return
			default:
				panic(errorValue)
			}
		}
	}()

	startThunk := func() packet {
		return evalNode(e, n)
	}

	return trampoline(startThunk), nil
}

func evalEachNode(e Env, ns []Node) []Node {
	result := make([]Node, len(ns))
	for i, n := range ns {
		evalNodeThunk := func() packet {
			return evalNode(e, n)
		}
		result[i] = trampoline(evalNodeThunk)
	}
	return result
}

func evalNode(e Env, n Node) packet {

	switch value := n.(type) {
	case *Number:
		return respond(value)
	case *Symbol:
		return respond(e.Get(value.Name))
	case *List:
		return bounce(func() packet { return evalList(e, value) })
	default:
		panicEvalError("Unknown form to evaluate: " + value.String())
	}

	return respond(&Symbol{Name: "nil"})
}

func evalDef(e Env, args []Node) packet {
	ensureArgCount("def", args, 2)

	name := toSymbolName(args[0])
	e.Set(name, trampoline(func() packet {
		return evalNode(e, args[1])
	}))
	return respond(&Symbol{Name: "nil"})
}

func evalEval(e Env, args []Node) packet {
	ensureArgCount("eval", args, 1)

	node := trampoline(func() packet {
		return evalNode(e, args[0])
	})

	return bounce(func() packet {
		return evalNode(e, node)
	})
}

func evalList(e Env, l *List) packet {
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
		case "apply":
			f := args[0]
			l := toListValue(trampoline(func() packet {
				return evalNode(e, args[1])
			}))
			nodes := append([]Node{f}, l.Nodes...)
			return respond(trampoline(func() packet {
				return evalList(e, &List{Nodes: nodes})
			}))
		case "def":
			return evalDef(e, args)
		case "eval":
			return evalEval(e, args)
		case "set!":
			name := toSymbolName(args[0])
			e.Update(name, trampoline(func() packet {
				return evalNode(e, args[1])
			}))
			return respond(&Symbol{Name: "nil"})
		case "if":
			predicate := toBooleanValue(trampoline(func() packet {
				return evalNode(e, args[0])
			}))
			if predicate {
				return bounce(func() packet {
					return evalNode(e, args[1])
				})
			}
			return bounce(func() packet {
				return evalNode(e, args[2])
			})
		case "cond":
			for i := 0; i < len(args); i += 2 {
				predicate := toBooleanValue(trampoline(func() packet {
					return evalNode(e, args[i])
				}))

				if predicate {
					return bounce(func() packet {
						return evalNode(e, args[i+1])
					})
				}
			}
			panicEvalError("No matching cond clause: " + l.String())
		case "fn":
			return bounce(func() packet {
				return evalFunctionDefinition(e, args)
			})
		case "quote":
			return respond(args[0])
		case "let":
			return bounce(func() packet {
				return evalLet(e, args)
			})
		}
	}

	headNode := trampoline(func() packet {
		return evalNode(e, head)
	})

	switch value := headNode.(type) {
	case *Primitive:
		f := value.Value
		return respond(f(e, evalEachNode(e, args)))
	case *Function:
		arguments := evalEachNode(e, args)

		return bounce(func() packet {
			return evalFunctionApplication(value, arguments)
		})
	default:
		panicEvalError("First item in list not a function: " + value.String())
	}

	return respond(&Symbol{Name: "nil"})
}

func evalFunctionApplication(f *Function, args []Node) packet {

	e := NewMapEnv(f.Name, f.ParentEnv)

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
		return evalNode(e, f.Body)
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

func evalLet(parentEnv Env, args []Node) packet {
	variableList := args[0]
	body := args[1]
	variableNodes := variableList.Children()

	e := NewMapEnv("let", parentEnv)

	// Evaluate variable assignments
	for i := 0; i < len(variableNodes); i += 2 {
		variable := variableNodes[i]
		expression := variableNodes[i+1]
		variableName := toSymbolName(variable)

		e.Set(variableName, trampoline(func() packet {
			return evalNode(e, expression)
		}))
	}

	// Evaluate body
	return bounce(func() packet {
		return evalNode(e, body)
	})
}

func ensureArgCount(formName string, args []Node, count int) {
	if len(args) != count {
		panicEvalError("Form '" + formName + "' requires exactly " + strconv.Itoa(count) + " argument(s), but was given " + strconv.Itoa(len(args)))
	}
}

func panicEvalError(s string) {
	panic(NewEvalError(s))
}

func toSymbolName(n Node) string {
	switch value := n.(type) {
	case *Symbol:
		return value.Name
	}

	panic("Not a symbol: " + n.String())
}
