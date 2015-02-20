package interpreter

import "vamos/lang/ast"

func evalSpecialLet(parentEnv Env, head ast.Node, args []ast.Node) packet {
	ensureSpecialArgsCountEquals("let", head, args, 2)

	body := args[1]

	var variableNodes ast.Nodes
	switch val := args[0].(type) {
	case ast.Coll:
		variableNodes = val.Children()
	default:
		panicEvalError(head, "Expected list as first argument to 'let': "+val.String())
	}

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

func evalSpecialGo(e Env, head ast.Node, args []ast.Node) packet {
	go evalEachNode(e, args)
	return respond(&ast.Nil{})
}

func evalSpecialBegin(e Env, head ast.Node, args []ast.Node) packet {
	results := evalEachNode(e, args)

	if len(results) == 0 {
		return respond(&ast.Nil{})
	}

	return respond(results[len(results)-1])
}

func evalSpecialDef(e Env, head ast.Node, args []ast.Node) packet {
	ensureSpecialArgsCountEquals("def", head, args, 2)

	name := toSymbolName(args[0])
	e.Set(name, trampoline(func() packet {
		return evalNode(e, args[1])
	}))
	return respond(&ast.Nil{})
}

func evalSpecialEval(e Env, head ast.Node, args []ast.Node) packet {
	ensureSpecialArgsCountInRange("eval", head, args, 1, 2)

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
			panicEvalError(args[0], "Second arg to 'eval' must be an environment: "+environmentNode.String())
			return respond(nil)
		}
	default:
		panicEvalError(args[0], "Unexpected number of args")
		return respond(nil)
	}
}

func evalSpecialFn(e Env, head ast.Node, args []ast.Node) packet {
	ensureSpecialArgsCountEquals("fn", head, args, 2)

	var parameterNodes ast.Nodes
	switch val := args[0].(type) {
	case ast.Coll:
		parameterNodes = val.Children()
	default:
		panicEvalError(head, "Expected list as first argument to 'fn': "+val.String())
	}

	return respond(&Function{
		Name:       "anonymous",
		Parameters: parameterNodes,
		Body:       args[1],
		ParentEnv:  e,
	})
}

func evalSpecialMacro(e Env, head ast.Node, args []ast.Node) packet {
	ensureSpecialArgsCountEquals("macro", head, args, 1)

	functionNode := trampoline(func() packet {
		return evalNode(e, args[0])
	})

	switch val := functionNode.(type) {
	case *Function:
		val.IsMacro = true
		return respond(val)
	default:
		panicEvalError(args[0], "macro expects a function argument but got: "+args[0].String())
		return respond(nil)
	}
}

func evalSpecialMacroexpand1(e Env, head ast.Node, args []ast.Node) packet {
	ensureSpecialArgsCountEquals("macroexpand1", head, args, 1)

	expansionNode := trampoline(func() packet {
		return evalNode(e, args[0])
	})

	switch value := expansionNode.(type) {
	case *ast.List:
		expansionResult := trampoline(func() packet {
			return evalList(e, value, false)
		})
		return respond(expansionResult)
	default:
		panicEvalError(args[0], "macroexpand1 expected a list but got: "+value.String())
		return respond(nil)
	}
}
