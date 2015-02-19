package lang

import (
	"fmt"
	"time"
	"vamos/util"
)

////////// Primitive Support

var trueSymbol, falseSymbol *Symbol

func initializePrimitives(e Env) {
	// Math
	addPrimitive(e, "+", 2, primAdd)
	addPrimitive(e, "-", 2, primSubtract)
	addPrimitive(e, "*", 2, primMult)
	addPrimitive(e, "/", 2, primDiv)
	addPrimitive(e, "<", 2, primLt)
	addPrimitive(e, ">", 2, primGt)

	// Equality
	addPrimitive(e, "=", 2, primEquals)

	addPrimitiveWithArityRange(e, "list", 0, -1, primList)

	// Collections
	addPrimitive(e, "first", 1, primFirst)
	addPrimitive(e, "rest", 1, primRest)
	addPrimitive(e, "cons", 2, primCons)
	addPrimitiveWithArityRange(e, "concat", 0, -1, primConcat)

	// Environments and types
	addPrimitive(e, "current-environment", 0, primCurrentEnvironment)
	addPrimitive(e, "typeof", 1, primTypeof)

	// Metaprogramming
	addPrimitive(e, "function-params", 1, primFunctionParams)
	addPrimitive(e, "function-body", 1, primFunctionBody)
	addPrimitive(e, "function-environment", 1, primFunctionEnvironment)

	// IO
	addPrimitive(e, "println", 1, primPrintln)
	addPrimitive(e, "load", 1, primLoad)
	addPrimitive(e, "now", 0, primNow)
	addPrimitive(e, "sleep", 1, primSleep)

	// Predefined symbols

	trueSymbol = &Symbol{Name: "true"}
	e.Set("true", trueSymbol)

	falseSymbol = &Symbol{Name: "false"}
	e.Set("false", falseSymbol)
}

func addPrimitiveWithArityRange(e Env, name string, minArity int, maxArity int, f primitiveFunction) {
	e.Set(
		name,
		NewPrimitive(name, minArity, maxArity, primitiveFunction(f)))
}

func addPrimitive(e Env, name string, arity int, f primitiveFunction) {
	e.Set(
		name,
		NewPrimitive(name, arity, arity, primitiveFunction(f)))
}

////////// Primitives

func primAdd(e Env, head Node, args []Node) Node {
	result := toNumberValue(args[0]) + toNumberValue(args[1])
	return &Number{Value: result}
}

func primSubtract(e Env, head Node, args []Node) Node {
	result := toNumberValue(args[0]) - toNumberValue(args[1])
	return &Number{Value: result}
}

func primEquals(e Env, head Node, args []Node) Node {
	if args[0].Equals(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primLt(e Env, head Node, args []Node) Node {
	if toNumberValue(args[0]) < toNumberValue(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primGt(e Env, head Node, args []Node) Node {
	if toNumberValue(args[0]) > toNumberValue(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primDiv(e Env, head Node, args []Node) Node {
	result := toNumberValue(args[0]) / toNumberValue(args[1])
	return &Number{Value: result}
}

func primMult(e Env, head Node, args []Node) Node {
	result := toNumberValue(args[0]) * toNumberValue(args[1])
	return &Number{Value: result}
}

func primList(e Env, head Node, args []Node) Node {
	return &ListNode{Nodes: args}
}

func primCurrentEnvironment(e Env, head Node, args []Node) Node {
	return NewEnvNode(e)
}

func primFunctionParams(e Env, head Node, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return NewListNode(val.Parameters)
	default:
		panicEvalError(args[0], "Argument to 'function-params' not a function: "+arg.String())
	}

	return nil
}

func primFunctionBody(e Env, head Node, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return val.Body
	default:
		panicEvalError(args[0], "Argument to 'function-body' not a function: "+arg.String())
	}

	return nil
}

func primFunctionEnvironment(e Env, head Node, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return NewEnvNode(val.ParentEnv)
	default:
		panicEvalError(args[0], "Argument to 'function-environment' not a function: "+arg.String())
	}

	return nil
}

func primTypeof(e Env, head Node, args []Node) Node {
	arg := args[0]
	return &Symbol{Name: arg.TypeName()}
}

func primPrintln(e Env, head Node, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *StringNode:
		fmt.Fprintf(writer, "%v\n", val.Value)
		return &NilNode{}
	}

	panicEvalError(arg, "Argument to 'println' not a string: "+arg.String())
	return nil
}

func primFirst(e Env, head Node, args []Node) Node {
	arg := args[0]

	switch val := arg.(type) {
	case Coll:
		return val.First()
	}

	panicEvalError(arg, "Cannot get first from a non-collection: "+arg.String())
	return nil
}

func primRest(e Env, head Node, args []Node) Node {
	arg := args[0]

	switch val := arg.(type) {
	case Coll:
		return val.Rest()
	}

	panicEvalError(arg, "Cannot get rest from a non-collection: "+arg.String())
	return nil
}

func primCons(e Env, head Node, args []Node) Node {
	sourceElement := args[0]
	targetColl := args[1]

	switch val := targetColl.(type) {
	case Coll:
		return val.Cons(sourceElement)
	}

	panicEvalError(sourceElement, "Cannot cons onto a non-collection: "+targetColl.String())
	return nil
}

func primConcat(e Env, head Node, args []Node) Node {
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

func primLoad(e Env, head Node, args []Node) Node {
	arg := args[0]
	switch val := arg.(type) {
	case *StringNode:
		fileName := val.Value

		if len(fileName) > 0 {
			content, err := util.ReadFile(fileName)
			if err != nil {
				panicEvalError(
					arg,
					fmt.Sprintf("Error while loading file <%v>: %v\n", fileName, err.Error()))
			} else {
				ParseEvalPrint(e, content, fileName, false)
			}
		}

		return &NilNode{}
	}

	panicEvalError(arg, "Argument to 'load' not a string: "+arg.String())
	return nil
}

func primNow(e Env, head Node, args []Node) Node {

	t := time.Now()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()

	result := NewListNode([]Node{
		&Number{Value: float64(year)},
		&Number{Value: float64(month)},
		&Number{Value: float64(day)},
		&Number{Value: float64(hour)},
		&Number{Value: float64(minute)},
		&Number{Value: float64(second)},
	})

	return result
}

func primSleep(e Env, head Node, args []Node) Node {

	arg := args[0]

	switch val := arg.(type) {
	case *Number:
		time.Sleep(time.Duration(val.Value) * time.Millisecond)
		return &NilNode{}
	}

	panicEvalError(arg, "Argument to 'sleep' not a number: "+arg.String())
	return nil
}
