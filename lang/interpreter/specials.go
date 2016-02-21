/*
A special form is a language primitive which does not evaluate its arguments before executing.
*/

package interpreter

import "vamos/lang/ast"

func specialQuote(e Env, head ast.Node, args []ast.Node) packet {
	return respond(args[0])
}

func specialUpdateBang(e Env, head ast.Node, args []ast.Node) packet {
	name := toSymbolName(args[0])
	rightHandSide := trampoline(func() packet {
		return evalNode(e, args[1])
	})
	if ok := e.Update(name, rightHandSide); !ok {
		panicEvalError(head, "Cannot 'update!' an undefined name: "+name)
	}
	return respond(&ast.Nil{})
}

func specialIf(e Env, head ast.Node, args []ast.Node) packet {
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
}

func specialCond(e Env, head ast.Node, args []ast.Node) packet {
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

	panicEvalError(head, "No matching cond clause: "+head.String())
	return respond(&ast.Nil{})
}

func specialLet(parentEnv Env, head ast.Node, args []ast.Node) packet {
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

func specialGo(e Env, head ast.Node, args []ast.Node) packet {
	go evalEachNode(e, args)
	return respond(&ast.Nil{})
}

func specialBegin(e Env, head ast.Node, args []ast.Node) packet {
	results := evalEachNode(e, args)

	if len(results) == 0 {
		return respond(&ast.Nil{})
	}

	return respond(results[len(results)-1])
}

func specialDef(e Env, head ast.Node, args []ast.Node) packet {
	name := toSymbolName(args[0])

	rightHandSide := trampoline(func() packet {
		return evalNode(e, args[1])
	})

	switch val := rightHandSide.(type) {
	case *Procedure:
		// Give a name to the procedure, allowing for better error messages
		val.Name = name
	}

	if _, exists := e.Get(name); exists {
		panicEvalError(head, "Cannot redefine a name: "+name)
	} else {
		e.Set(name, rightHandSide)
	}

	return respond(&ast.Nil{})
}

func specialEval(e Env, head ast.Node, args []ast.Node) packet {
	checkSpecialArgs("eval", head, args, 1, 2)

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

func specialFn(e Env, head ast.Node, args []ast.Node) packet {

	var parameterNodes ast.Nodes
	switch val := args[0].(type) {
	case ast.Coll:
		parameterNodes = val.Children()
	default:
		panicEvalError(head, "Expected list as first argument to 'proc': "+val.String())
	}

	return respond(&Procedure{
		Name:       "anonymous",
		Parameters: parameterNodes,
		Body:       args[1],
		ParentEnv:  e,
	})
}

func specialMacro(e Env, head ast.Node, args []ast.Node) packet {

	procedureNode := trampoline(func() packet {
		return evalNode(e, args[0])
	})

	switch val := procedureNode.(type) {
	case *Procedure:
		val.IsMacro = true
		return respond(val)
	default:
		panicEvalError(args[0], "macro expects a procedure argument but got: "+args[0].String())
		return respond(nil)
	}
}

func specialMacroexpand1(e Env, head ast.Node, args []ast.Node) packet {

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
