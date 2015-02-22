package interpreter

import (
	"fmt"
	"time"
	"vamos/lang/ast"
	"vamos/lang/parser"
	"vamos/util"
)

////////// Primitive Support

var trueSymbol, falseSymbol *ast.Symbol

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
	addPrimitive(e, "read-string", 1, primReadString)

	// IO
	addPrimitive(e, "println", 1, primPrintln)
	addPrimitive(e, "load", 1, primLoad)
	addPrimitive(e, "now", 0, primNow)
	addPrimitive(e, "sleep", 1, primSleep)

	// Predefined symbols

	trueSymbol = &ast.Symbol{Name: "true"}
	e.Set("true", trueSymbol)

	falseSymbol = &ast.Symbol{Name: "false"}
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

func primAdd(e Env, head ast.Node, args []ast.Node) ast.Node {
	result := toNumberValue(args[0]) + toNumberValue(args[1])
	return &ast.Number{Value: result}
}

func primSubtract(e Env, head ast.Node, args []ast.Node) ast.Node {
	result := toNumberValue(args[0]) - toNumberValue(args[1])
	return &ast.Number{Value: result}
}

func primEquals(e Env, head ast.Node, args []ast.Node) ast.Node {
	if args[0].Equals(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primLt(e Env, head ast.Node, args []ast.Node) ast.Node {
	if toNumberValue(args[0]) < toNumberValue(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primGt(e Env, head ast.Node, args []ast.Node) ast.Node {
	if toNumberValue(args[0]) > toNumberValue(args[1]) {
		return trueSymbol
	}
	return falseSymbol
}

func primDiv(e Env, head ast.Node, args []ast.Node) ast.Node {
	result := toNumberValue(args[0]) / toNumberValue(args[1])
	return &ast.Number{Value: result}
}

func primMult(e Env, head ast.Node, args []ast.Node) ast.Node {
	result := toNumberValue(args[0]) * toNumberValue(args[1])
	return &ast.Number{Value: result}
}

func primList(e Env, head ast.Node, args []ast.Node) ast.Node {
	return &ast.List{Nodes: args}
}

func primCurrentEnvironment(e Env, head ast.Node, args []ast.Node) ast.Node {
	return NewEnvNode(e)
}

func primFunctionParams(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return ast.NewList(val.Parameters)
	default:
		panicEvalError(args[0], "Argument to 'function-params' not a function: "+arg.String())
	}

	return nil
}

func primFunctionBody(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return val.Body
	default:
		panicEvalError(args[0], "Argument to 'function-body' not a function: "+arg.String())
	}

	return nil
}

func primFunctionEnvironment(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Function:
		return NewEnvNode(val.ParentEnv)
	default:
		panicEvalError(args[0], "Argument to 'function-environment' not a function: "+arg.String())
	}

	return nil
}

func primTypeof(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	return &ast.Symbol{Name: arg.TypeName()}
}

func primPrintln(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	switch val := arg.(type) {
	case *ast.Str:
		fmt.Fprintf(writer, "%v\n", val.Value)
		return &ast.Nil{}
	}

	panicEvalError(arg, "Argument to 'println' not a string: "+arg.String())
	return nil
}

func primFirst(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]

	switch val := arg.(type) {
	case ast.Coll:
		return val.First()
	}

	panicEvalError(arg, "Cannot get first from a non-collection: "+arg.String())
	return nil
}

func primRest(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]

	switch val := arg.(type) {
	case ast.Coll:
		return val.Rest()
	}

	panicEvalError(arg, "Cannot get rest from a non-collection: "+arg.String())
	return nil
}

func primCons(e Env, head ast.Node, args []ast.Node) ast.Node {
	sourceElement := args[0]
	targetColl := args[1]

	switch val := targetColl.(type) {
	case ast.Coll:
		result, err := val.Cons(sourceElement)
		if err != nil {
			panicEvalError(head, err.Error())
			return nil
		}
		return result
	}

	panicEvalError(sourceElement, "Cannot cons onto a non-collection: "+targetColl.String())
	return nil
}

func primConcat(e Env, head ast.Node, args []ast.Node) ast.Node {
	var sum ast.Node

	for _, arg := range args {
		if sum == nil {
			sum = arg
		} else {
			switch sumVal := sum.(type) {
			case ast.Coll:
				switch argVal := arg.(type) {
				case ast.Coll:
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
		return &ast.Nil{}
	} else {
		return sum
	}
}

func primLoad(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	switch val := arg.(type) {
	case *ast.Str:
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

		return &ast.Nil{}
	}

	panicEvalError(arg, "Argument to 'load' not a string: "+arg.String())
	return nil
}

func primNow(e Env, head ast.Node, args []ast.Node) ast.Node {

	t := time.Now()
	year, month, day := t.Date()
	hour, minute, second := t.Clock()

	result := ast.NewList([]ast.Node{
		&ast.Number{Value: float64(year)},
		&ast.Number{Value: float64(month)},
		&ast.Number{Value: float64(day)},
		&ast.Number{Value: float64(hour)},
		&ast.Number{Value: float64(minute)},
		&ast.Number{Value: float64(second)},
	})

	return result
}

func primSleep(e Env, head ast.Node, args []ast.Node) ast.Node {

	arg := args[0]

	switch val := arg.(type) {
	case *ast.Number:
		time.Sleep(time.Duration(val.Value) * time.Millisecond)
		return &ast.Nil{}
	}

	panicEvalError(arg, "Argument to 'sleep' not a number: "+arg.String())
	return nil
}

func primReadString(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	switch val := arg.(type) {
	case *ast.Str:
		nodes, parseErrors := parser.Parse(val.Value, "string")

		if parseErrors != nil {
			panicEvalError(arg, fmt.Sprintf("Unable to read string %v: %v", val, parseErrors))
			return nil
		}

		if len(nodes) == 0 {
			return &ast.Nil{}
		}

		return nodes[0]
	}

	panicEvalError(arg, "Argument to 'read-string' not a string: "+arg.String())
	return nil
}
