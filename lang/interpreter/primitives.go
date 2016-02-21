package interpreter

import (
	"bytes"
	"fmt"
	"strings"
	"time"
	"vamos/lang/ast"
	"vamos/lang/parser"
	"vamos/util"
)

////////// Primitive Support

var trueSymbol, falseSymbol *ast.Symbol

func initializePrimitives(e Env) {
	// Basic
	addPrimitiveWithArityRange(e, "list", 0, -1, primList)
	addPrimitive(e, "apply", 2, primApply)

	// Math
	addPrimitive(e, "+", 2, primAdd)
	addPrimitive(e, "-", 2, primSubtract)
	addPrimitive(e, "*", 2, primMult)
	addPrimitive(e, "/", 2, primDiv)
	addPrimitive(e, "<", 2, primLt)
	addPrimitive(e, ">", 2, primGt)

	// Strings
	addPrimitiveWithArityRange(e, "str", 0, -1, primStr)

	// Equality
	addPrimitive(e, "=", 2, primEquals)

	// Collections
	addPrimitive(e, "first", 1, primFirst)
	addPrimitive(e, "rest", 1, primRest)
	addPrimitive(e, "cons", 2, primCons)
	addPrimitiveWithArityRange(e, "concat", 0, -1, primConcat)

	// Environments and types
	addPrimitive(e, "current-environment", 0, primCurrentEnvironment)
	addPrimitive(e, "typeof", 1, primTypeof)
	addPrimitive(e, "update-element!", 3, primUpdateElementBang)

	// Metaprogramming
	addPrimitive(e, "routine-params", 1, primProcedureParams)
	addPrimitive(e, "routine-body", 1, primProcedureBody)
	addPrimitive(e, "routine-environment", 1, primProcedureEnvironment)
	addPrimitive(e, "read-string", 1, primReadString)
	addPrimitive(e, "readable-string", 1, primReadableString)

	// IO
	addPrimitiveWithArityRange(e, "println", 1, -1, primPrintln)
	addPrimitiveWithArityRange(e, "read-line", 0, 0, primReadLine)
	addPrimitive(e, "load", 1, primLoad)
	addPrimitive(e, "now", 0, primNow)
	addPrimitive(e, "sleep", 1, primSleep)
	addPrimitiveWithArityRange(e, "panic", 0, -1, primPanic)

	// Concurrency
	addPrimitive(e, "chan", 0, primChan)
	addPrimitive(e, "send!", 2, primSendBang)
	addPrimitive(e, "take!", 1, primTakeBang)
	addPrimitive(e, "close!", 1, primCloseBang)

	// Predefined symbols

	trueSymbol = &ast.Symbol{Name: "true"}
	e.Set("true", trueSymbol)

	falseSymbol = &ast.Symbol{Name: "false"}
	e.Set("false", falseSymbol)
}

func addPrimitiveWithArityRange(e Env, name string, minArity int, maxArity int, f primitiveFunc) {
	e.Set(
		name,
		NewPrimitive(name, minArity, maxArity, primitiveFunc(f)))
}

func addPrimitive(e Env, name string, arity int, f primitiveFunc) {
	e.Set(
		name,
		NewPrimitive(name, arity, arity, primitiveFunc(f)))
}

////////// Primitives

func primApply(e Env, head ast.Node, args []ast.Node) ast.Node {
	evaluatedHead := args[0]
	switch headVal := evaluatedHead.(type) {
	case Routine:
		l := toListValue(args[1])
		return trampoline(func() packet {
			return evalInvokeRoutine(e, headVal, head, l.Nodes, true)
		})
	default:
		panicEvalError(head, "First argument to 'apply' not a routine: "+headVal.String())
		return &ast.Nil{}
	}
}

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

func primProcedureParams(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Procedure:
		return ast.NewList(val.Parameters)
	default:
		panicEvalError(args[0], "Argument to 'routine-params' not a procedure: "+arg.String())
	}

	return nil
}

func primProcedureBody(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Procedure:
		return val.Body
	default:
		panicEvalError(args[0], "Argument to 'routine-body' not a procedure: "+arg.String())
	}

	return nil
}

func primProcedureEnvironment(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	switch val := arg.(type) {
	case *Procedure:
		return NewEnvNode(val.ParentEnv)
	default:
		panicEvalError(args[0], "Argument to 'routine-environment' not a procedure: "+arg.String())
	}

	return nil
}

func primTypeof(e Env, head ast.Node, args []ast.Node) ast.Node {
	arg := args[0]
	return &ast.Symbol{Name: arg.TypeName()}
}

func primPanic(e Env, head ast.Node, args []ast.Node) ast.Node {
	var buffer bytes.Buffer

	for i, arg := range args {
		if i > 0 {
			buffer.WriteString(" ")
		}

		switch val := arg.(type) {
		case *ast.Str:
			buffer.WriteString(val.Value)
		case ast.Node:
			buffer.WriteString(val.String())
		default:
			panicEvalError(arg, "Unrecognized argument type to 'panic': "+arg.String())
		}
	}

	panicApplicationError(head, buffer.String())
	return &ast.Nil{}
}

func primPrintln(e Env, head ast.Node, args []ast.Node) ast.Node {
	for i, arg := range args {
		if i > 0 {
			fmt.Fprintf(writer, " ")
		}

		switch val := arg.(type) {
		case *ast.Str:
			fmt.Fprintf(writer, "%v", val.Value)
		case ast.Node:
			fmt.Fprintf(writer, "%v", val.String())
		default:
			fmt.Fprintf(writer, "\n")
			panicEvalError(arg, "Unrecognized argument type to 'println': "+arg.String())
		}
	}

	fmt.Fprintf(writer, "\n")
	return &ast.Nil{}
}

func primReadLine(e Env, head ast.Node, args []ast.Node) ast.Node {
	s := readLine() // TODO: uses a global variable :(
	trimmed := strings.TrimSuffix(s, "\n")
	return ast.NewStr(trimmed)
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

func primUpdateElementBang(e Env, head ast.Node, args []ast.Node) ast.Node {
	leftHandSide := args[0]

	indexNode := args[1]
	indexNumber, ok := indexNode.(*ast.Number)
	if !ok {
		panicEvalError(head, "Index in 'update-element!' is not a number: "+indexNode.String())
	}
	index := int(indexNumber.Value)

	rightHandSide := args[2]

	switch val := leftHandSide.(type) {
	case *ast.List:
		children := val.Children()
		children[index] = rightHandSide
	default:
		panicEvalError(head, "Cannot 'update-element!' in a non-list: "+leftHandSide.String())
	}

	return &ast.Nil{}
}

func primReadableString(e Env, head ast.Node, args []ast.Node) ast.Node {
	return ast.NewStr(args[0].String())
}

func primStr(e Env, head ast.Node, args []ast.Node) ast.Node {
	var buffer bytes.Buffer

	// TODO replace below with calls to .FriendlyString()
	for _, arg := range args {
		switch val := arg.(type) {
		case *ast.Str:
			buffer.WriteString(val.Value)
		case *ast.Char:
			buffer.WriteRune(val.Value)
		case ast.Node:
			buffer.WriteString(val.String())
		default:
			panicEvalError(arg, "Unrecognized argument type to 'str': "+arg.String())
		}
	}

	return ast.NewStr(buffer.String())
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
				ParseEvalPrint(e, content, readLine, fileName, false)
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

func primChan(e Env, head ast.Node, args []ast.Node) ast.Node {
	return NewChan()
}

func primSendBang(e Env, head ast.Node, args []ast.Node) ast.Node {
	chanArg := args[0]
	switch chanVal := chanArg.(type) {
	case *Chan:
		messageArg := args[1]
		chanVal.Value <- messageArg
	default:
		panicEvalError(head, "Target of a send! must be a chan: "+chanArg.String())
	}

	return &ast.Nil{}
}

func primTakeBang(e Env, head ast.Node, args []ast.Node) ast.Node {
	chanArg := args[0]
	switch chanVal := chanArg.(type) {
	case *Chan:
		n, more := <-chanVal.Value
		if !more {
			return &ast.Nil{}
		}
		return n
	default:
		panicEvalError(head, "Source of a take! must be a chan: "+chanArg.String())
	}

	return &ast.Nil{}
}

func primCloseBang(e Env, head ast.Node, args []ast.Node) ast.Node {
	chanArg := args[0]
	switch chanVal := chanArg.(type) {
	case *Chan:
		close(chanVal.Value)
	default:
		panicEvalError(head, "Argument to 'close!' must be a chan: "+chanArg.String())
	}

	return &ast.Nil{}
}
