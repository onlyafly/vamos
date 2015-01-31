package interpretation

import . "vamos/lang/ast"

func evalSpecialDef(e Env, args []Node) packet {
	ensureArgCount("def", args, 2)

	name := toSymbolName(args[0])
	e.Set(name, trampoline(func() packet {
		return evalNode(e, args[1])
	}))
	return respond(&Symbol{Name: "nil"})
}

func evalSpecialEval(e Env, args []Node) packet {
	ensureArgCount("eval", args, 1)

	node := trampoline(func() packet {
		return evalNode(e, args[0])
	})

	return bounce(func() packet {
		return evalNode(e, node)
	})
}

func evalSpecialFn(e Env, args []Node) packet {
	ensureArgCount("fn", args, 2)

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
	ensureArgCount("macro", args, 2)

	parameterList := args[0]
	parameterNodes := parameterList.Children()

	return respond(&Macro{
		Name:       "anonymous",
		Parameters: parameterNodes,
		Body:       args[1],
		ParentEnv:  e,
	})
}
