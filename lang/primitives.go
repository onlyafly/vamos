package lang

////////// Primitive Support

type primitiveFunction func(Env, []Node) Node

var trueSymbol, falseSymbol *Symbol

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

	addPrimitive(e, "env", primInspectEnv)

	trueSymbol = &Symbol{Name: "true"}
	falseSymbol = &Symbol{Name: "false"}
	e.Set("true", trueSymbol)
	e.Set("false", falseSymbol)
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

func primInspectEnv(e Env, args []Node) Node {
	print(e.String(), "\n")
	return falseSymbol
}
