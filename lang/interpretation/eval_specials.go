package interpretation

import . "vamos/lang/ast"

func evalSpecialLet(parentEnv Env, args []Node) packet {
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

func evalSpecialBegin(e Env, args []Node) packet {
	results := evalEachNode(e, args)

	if len(results) == 0 {
		return respond(&Symbol{Name: "nil"})
	}

	return respond(results[len(results)-1])
}

func evalSpecialDef(e Env, args []Node) packet {
	ensureSpecialArgsCountEquals("def", args, 2)

	name := toSymbolName(args[0])
	e.Set(name, trampoline(func() packet {
		return evalNode(e, args[1])
	}))
	return respond(&Symbol{Name: "nil"})
}

func evalSpecialEval(e Env, args []Node) packet {
	ensureSpecialArgsCountInRange("eval", args, 1, 2)

	node := trampoline(func() packet {
		return evalNode(e, args[0])
	})

	switch len(args) {
	case 1:
		return bounce(func() packet {
			return evalNode(e, node)
		})
	case 2:
		nodeArg1 := trampoline(func() packet {
			return evalNode(e, args[1])
		})

		switch environmentNode := nodeArg1.(type) {
		case *EnvNode:
			return bounce(func() packet {
				return evalNode(environmentNode.Env, node)
			})
		default:
			panicEvalError("Second arg to 'eval' must be an environment: " + environmentNode.String())
			return respond(nil)
		}
	default:
		panicEvalError("Unexpected number of args")
		return respond(nil)
	}
}

func evalSpecialFn(e Env, args []Node) packet {
	ensureSpecialArgsCountEquals("fn", args, 2)

	parameterList := args[0]
	parameterNodes := parameterList.Children()

	return respond(&Function{
		Name:       "anonymous",
		Parameters: parameterNodes,
		Body:       args[1],
		ParentEnv:  e,
	})
}

func evalSpecialMacro(e Env, args []Node) packet {
	ensureSpecialArgsCountEquals("macro", args, 2)

	parameterList := args[0]
	parameterNodes := parameterList.Children()

	return respond(&Macro{
		Name:       "anonymous",
		Parameters: parameterNodes,
		Body:       args[1],
		ParentEnv:  e,
	})
}

func evalSpecialMacroexpand1(e Env, args []Node) packet {
	ensureSpecialArgsCountEquals("macroexpand1", args, 1)

	expansionNode := trampoline(func() packet {
		return evalNode(e, args[0])
	})

	switch value := expansionNode.(type) {
	case *List:
		expansionResult := trampoline(func() packet {
			return evalList(e, value, false)
		})
		return respond(expansionResult)
	default:
		panicEvalError("macroexpand1 expected a list but got: " + value.String())
		return respond(nil)
	}
}
