package lang

import (
	"fmt"
)

////////// Primitive Support

var trueSymbol, falseSymbol *Symbol

func initializePrimitives(e Env) {
	// Math
	addPrimitive(e, "+", primAdd)
	addPrimitive(e, "-", primSubtract)
	addPrimitive(e, "*", primMult)
	addPrimitive(e, "/", primDiv)
	addPrimitive(e, "<", primLt)
	addPrimitive(e, ">", primGt)

	// Equality
	addPrimitive(e, "=", primEquals)

	addPrimitive(e, "list", primList)

	// Collections
	addPrimitive(e, "first", primFirst)
	addPrimitive(e, "rest", primRest)
	addPrimitive(e, "cons", primCons)
	addPrimitive(e, "concat", primConcat)

	// Environments and types
	addPrimitive(e, "current-environment", primCurrentEnvironment)
	addPrimitive(e, "typeof", primTypeof)

	// Metaprogramming
	addPrimitive(e, "function-params", primFunctionParams)
	addPrimitive(e, "function-body", primFunctionBody)
	addPrimitive(e, "function-environment", primFunctionEnvironment)

	// IO
	addPrimitive(e, "println", primPrintln)

	// Predefined symbols

	trueSymbol = &Symbol{Name: "true"}
	e.Set("true", trueSymbol)

	falseSymbol = &Symbol{Name: "false"}
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
	return &ListNode{Nodes: args}
}

func primFirst(e Env, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *ListNode:
		return val.Nodes[0]
	}

	panicEvalError(args[0], "Argument to first not a list: "+arg.String())
	return nil
}

func primRest(e Env, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *ListNode:
		return &ListNode{Nodes: val.Nodes[1:]}
	}

	panicEvalError(args[0], "Argument to rest not a list: "+arg.String())
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
		panicEvalError(args[0], "Argument to 'function-params' not a function: "+arg.String())
	}

	return nil
}

func primFunctionBody(e Env, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return val.Body
	default:
		panicEvalError(args[0], "Argument to 'function-body' not a function: "+arg.String())
	}

	return nil
}

func primFunctionEnvironment(e Env, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return NewEnvNode(val.ParentEnv)
	default:
		panicEvalError(args[0], "Argument to 'function-environment' not a function: "+arg.String())
	}

	return nil
}

func primTypeof(e Env, args []Node) Node {
	arg := args[0]
	return &Symbol{Name: arg.TypeName()}
}

func primPrintln(e Env, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *StringNode:
		fmt.Fprintf(writer, "%v\n", val.Value)
		return &NilNode{}
	}

	panicEvalError(arg, "Argument to 'println' not a string: "+arg.String())
	return nil
}

func primCons(e Env, args []Node) Node {
	sourceElement := args[0]
	targetColl := args[1]

	switch val := targetColl.(type) {
	case Coll:
		return val.Cons(sourceElement)
	}

	panicEvalError(sourceElement, "Cannot cons onto a non-collection: "+targetColl.String())
	return nil
}

func primConcat(e Env, args []Node) Node {
	var sum Node = nil

	for _, arg := range args {
		if sum == nil {
			sum = arg
		} else {
			switch sumVal := sum.(type) {
			case Coll:
				switch argVal := arg.(type) {
				case Coll:
					sum = sumVal.Append(argVal)
				default:
					panicEvalError(arg, "Cannot concat a collection with a non-collection: "+arg.String())
				}
			default:
				panicEvalError(arg, "Cannot concat a non-collection type: "+sum.String())
			}
		}
	}

	if sum == nil {
		return &NilNode{}
	} else {
		return sum
	}
}
