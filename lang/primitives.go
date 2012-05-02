package lang

////////// Primitive Support

type primitiveFunction func([]Node) Node

var trueSymbol, falseSymbol *Symbol

func initializePrimitives(e Env) {
	addPrimitive(e, "+", primAdd)
	addPrimitive(e, "-", primSubtract)
	addPrimitive(e, "*", primMult)
	addPrimitive(e, "/", primDiv)
	addPrimitive(e, "=", primEquals)
	addPrimitive(e, "<", primLt)
	addPrimitive(e, ">", primGt)

	trueSymbol = &Symbol{Name: "true"}
	falseSymbol = &Symbol{Name: "false"}
	e.Set("true", trueSymbol)
	e.Set("false", falseSymbol)
}

func addPrimitive(e Env, name string, f func([]Node) Node) {
	e.Set(
		name,
		&Primitive{Name: name, Value: primitiveFunction(f)})
}

////////// Primitives

func primAdd(args []Node) Node {
	result := toNumberValue(args[0]) + toNumberValue(args[1])
	return &Number{Value: result}
}

func primSubtract(args []Node) Node {
	result := toNumberValue(args[0]) - toNumberValue(args[1])
	return &Number{Value: result}
}

func primEquals(args []Node) Node {
	if toNumberValue(args[0]) == toNumberValue(args[1]) {
		return trueSymbol
	}

	return falseSymbol
}

func primLt(args []Node) Node {
	if toNumberValue(args[0]) < toNumberValue(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primGt(args []Node) Node {
	if toNumberValue(args[0]) > toNumberValue(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primDiv(args []Node) Node {
	result := toNumberValue(args[0]) / toNumberValue(args[1])
	return &Number{Value: result}
}

func primMult(args []Node) Node {
	result := toNumberValue(args[0]) * toNumberValue(args[1])
	return &Number{Value: result}
}
