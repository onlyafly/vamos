package interpretation

import . "vamos/lang/ast"

////////// Primitive Support

var trueSymbol, falseSymbol, nilSymbol *Symbol

func initializePrimitives(e Env) {
	addPrimitive(e, "+", primAdd)
	addPrimitive(e, "-", primSubtract)
	addPrimitive(e, "*", primMult)
	addPrimitive(e, "/", primDiv)
	addPrimitive(e, "=", primEquals)
	addPrimitive(e, "<", primLt)
	addPrimitive(e, ">", primGt)
	addPrimitive(e, "list", primList)
	addPrimitive(e, "first", primFirst)
	addPrimitive(e, "rest", primRest)

	addPrimitive(e, "current-environment", primCurrentEnvironment)
	addPrimitive(e, "typeof", primTypeof)

	addPrimitive(e, "function-params", primFunctionParams)
	addPrimitive(e, "function-body", primFunctionBody)
	addPrimitive(e, "function-environment", primFunctionEnvironment)

	trueSymbol = &Symbol{Name: "true"}
	falseSymbol = &Symbol{Name: "false"}
	e.Set("true", trueSymbol)
	e.Set("false", falseSymbol)

	nilSymbol = &Symbol{Name: "nil"}
	e.Set("nil", nilSymbol)
}

func addPrimitive(e Env, name string, f primitiveFunction) {
	e.Set(
		name,
		&Primitive{Name: name, Value: primitiveFunction(f)})
}

////////// Primitives

func primAdd(e Env, args []Node) Node {
	result := toNumberValue(args[0]) + toNumberValue(args[1])
	return &Number{Value: result}
}

func primSubtract(e Env, args []Node) Node {
	result := toNumberValue(args[0]) - toNumberValue(args[1])
	return &Number{Value: result}
}

func primEquals(e Env, args []Node) Node {
	if args[0].Equals(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primLt(e Env, args []Node) Node {
	if toNumberValue(args[0]) < toNumberValue(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primGt(e Env, args []Node) Node {
	if toNumberValue(args[0]) > toNumberValue(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primDiv(e Env, args []Node) Node {
	result := toNumberValue(args[0]) / toNumberValue(args[1])
	return &Number{Value: result}
}

func primMult(e Env, args []Node) Node {
	result := toNumberValue(args[0]) * toNumberValue(args[1])
	return &Number{Value: result}
}

func primList(e Env, args []Node) Node {
	return &List{Nodes: args}
}

func primFirst(e Env, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *List:
		return val.Nodes[0]
	}

	panicEvalError("Argument to first not a list: " + arg.String())
	return nil
}

func primRest(e Env, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *List:
		return &List{Nodes: val.Nodes[1:]}
	}

	panicEvalError("Argument to rest not a list: " + arg.String())
	return nil
}

func primCurrentEnvironment(e Env, args []Node) Node {
	return NewEnvNode(e)
}

func primFunctionParams(e Env, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return NewList(val.Parameters)
	default:
		panicEvalError("Argument to 'function-params' not a function: " + arg.String())
	}

	return nil
}

func primFunctionBody(e Env, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return val.Body
	default:
		panicEvalError("Argument to 'function-body' not a function: " + arg.String())
	}

	return nil
}

func primFunctionEnvironment(e Env, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return NewEnvNode(val.ParentEnv)
	default:
		panicEvalError("Argument to 'function-environment' not a function: " + arg.String())
	}

	return nil
}

func primTypeof(e Env, args []Node) Node {
	arg := args[0]
	return &Symbol{Name: arg.TypeName()}
}
