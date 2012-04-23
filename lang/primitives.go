package lang

////////// Primitive Support

type primitiveFunction func([]Node) Node

var trueSymbol, falseSymbol *Symbol

func initializePrimitives(e Env) {
	addPrimitive(e, "+", primAdd)
	addPrimitive(e, "-", primSubtract)
	addPrimitive(e, "=", primEquals)

	trueSymbol = &Symbol{Name: "true"}
	falseSymbol = &Symbol{Name: "false"}
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
