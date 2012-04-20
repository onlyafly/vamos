package lang

////////// Primitive Support

type primitiveFunction func([]Node) Node

func initializePrimitives(e Env) {
	addPrimitive(e, "+", primAdd)
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
