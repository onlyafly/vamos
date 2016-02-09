package interpreter

import (
	"fmt"
	"io"
	"vamos/lang/ast"
)

// TODO get rid of this global variable
var writer io.Writer
var readLine func() string

// Eval evaluates a node in an environment.
func Eval(e Env, n ast.Node, w io.Writer, rl func() string) (result ast.Node, err error) {
	defer func() {
		if e := recover(); e != nil {
			result = nil
			switch errorValue := e.(type) {
			case *EvalError:
				err = errorValue
				return
			default:
				panic(errorValue)
			}
		}
	}()

	writer = w
	readLine = rl

	startThunk := func() packet {
		return evalNode(e, n)
	}

	return trampoline(startThunk), nil
}

func evalEachNode(e Env, ns []ast.Node) []ast.Node {
	result := make([]ast.Node, len(ns))
	for i, n := range ns {
		evalNodeThunk := func() packet {
			return evalNode(e, n)
		}
		result[i] = trampoline(evalNodeThunk)
	}
	return result
}

func evalNode(e Env, n ast.Node) packet {

	switch value := n.(type) {
	case *ast.Number:
		return respond(value)
	case *ast.Symbol:
		result, ok := e.Get(value.Name)
		if !ok {
			panicEvalError(value, "Name not defined: "+value.Name)
		}
		return respond(result)
	case *ast.Str:
		return respond(value)
	case *ast.Char:
		return respond(value)
	case *ast.List:
		return bounce(func() packet { return evalList(e, value, true) })
	case *ast.Nil:
		return respond(&ast.Nil{})
	default:
		panicEvalError(n, "Unknown form to evaluate: "+value.String())
	}

	return respond(&ast.Nil{})
}

func evalList(e Env, l *ast.List, shouldEvalMacros bool) packet {
	elements := l.Nodes

	if len(elements) == 0 {
		panicEvalError(l, "Empty list cannot be evaluated: "+l.String())
		return respond(nil)
	}

	/*
		Ten Primitives

		McCarthy introduced the ten primitives of lisp in 1960. All other pure lisp
		functions (i.e. all functions which don't do I/O or interact with the environment)
		can be implemented with these primitives. Thus, when implementing or porting
		lisp, these are the only functions which need to be implemented in a lower
		language. The way the non-primitives of lisp can be constructed from primitives
		is analogous to the way theorems can be proven from axioms in mathematics.

		The primitives are:

		Lisp:  atom  quote eq car   cdr  cons cond lambda label apply
		Vamos: atom? quote =  first rest cons cond fn			def		apply
	*/

	head := elements[0]
	args := elements[1:]

	switch value := head.(type) {
	case *ast.Symbol:
		switch value.Name {
		case "apply":
			checkSpecialArgs("apply", head, args, 2, 2)
			return specialApply(e, head, args)
		case "def":
			checkSpecialArgs("def", head, args, 2, 2)
			return specialDef(e, head, args)
		case "eval":
			checkSpecialArgs("eval", head, args, 1, 2)
			return specialEval(e, head, args)
		case "update!":
			checkSpecialArgs("update!", head, args, 2, 2)
			return specialUpdateBang(e, head, args)
		/*case "update-element!":
		checkSpecialArgs("update-element!", head, args, 3, 3)
		return specialUpdateElementBang(e, head, args)*/
		case "if":
			checkSpecialArgs("if", head, args, 3, 3)
			return specialIf(e, head, args)
		case "cond":
			checkSpecialArgs("cond", head, args, 2, -1)
			return specialCond(e, head, args)
		case "fn":
			checkSpecialArgs("fn", head, args, 2, 2)
			return specialFn(e, head, args)
		case "macro":
			checkSpecialArgs("macro", head, args, 1, 1)
			return specialMacro(e, head, args)
		case "macroexpand1":
			checkSpecialArgs("macroexpand1", head, args, 1, 1)
			return specialMacroexpand1(e, head, args)
		case "quote":
			checkSpecialArgs("quote", head, args, 1, 1)
			return specialQuote(e, head, args)
		case "let":
			checkSpecialArgs("let", head, args, 2, 2)
			return specialLet(e, head, args)
		case "begin":
			checkSpecialArgs("begin", head, args, 0, -1)
			return specialBegin(e, head, args)
		case "go":
			checkSpecialArgs("go", head, args, 0, -1)
			return specialGo(e, head, args)
		}
	}

	headNode := trampoline(func() packet {
		return evalNode(e, head)
	})

	switch value := headNode.(type) {
	case *Primitive:
		f := value.Value
		checkPrimitiveArgs(value.Name, head, args, value.MinArity, value.MaxArity)
		return respond(f(e, head, evalEachNode(e, args)))
	case *Function:
		return bounce(func() packet {
			return evalFunctionApplication(e, value, head, args, shouldEvalMacros)
		})
	default:
		panicEvalError(head, "First item in list not a function: "+value.String())
	}

	return respond(&ast.Nil{})
}

func evalFunctionApplication(dynamicEnv Env, f *Function, head ast.Node, unevaledArgs ast.Nodes, shouldEvalMacros bool) packet {
	defer func() {
		if e := recover(); e != nil {
			switch errorValue := e.(type) {
			case *EvalError:
				fmt.Printf("TRACE: (%v: %v): call to %v\n", head.Loc().Filename, head.Loc().Line, f.Name)
				panic(errorValue)
			default:
				panic(errorValue)
			}
		}
	}()

	// Validate parameters
	isVariableNumberOfParams := false
	for _, param := range f.Parameters {
		switch paramVal := param.(type) {
		case *ast.Symbol:
			if paramVal.Name == "&rest" {
				isVariableNumberOfParams = true
			}
		default:
			panicEvalError(head, "Function parameters should only be symbols: "+param.String())
		}
	}
	if !isVariableNumberOfParams {
		if len(unevaledArgs) != len(f.Parameters) {
			panicEvalError(head, fmt.Sprintf(
				"Function '%v' expects %v argument(s), but was given %v. Function parameter list: %v. Arguments: %v.",
				f.Name,
				len(f.Parameters),
				len(unevaledArgs),
				f.Parameters,
				unevaledArgs))
		}
	}

	// Create the lexical environment based on the function's lexical parent
	lexicalEnv := NewMapEnv(f.Name, f.ParentEnv)

	// Prepare the arguments for application
	var args []ast.Node
	if f.IsMacro {
		args = unevaledArgs
	} else {
		args = evalEachNode(dynamicEnv, unevaledArgs)
	}

	// Map arguments to parameters
	isMappingRestArgs := false
	iarg := 0
	for iparam, param := range f.Parameters {
		paramName := toSymbolName(param)
		if isMappingRestArgs {
			restArgs := args[iarg:]
			restList := ast.NewList(restArgs)
			lexicalEnv.Set(paramName, restList)
		} else if paramName == "&rest" {
			isMappingRestArgs = true
		} else {
			lexicalEnv.Set(paramName, args[iparam])
			iarg++
		}
	}

	if f.IsMacro {
		expandedMacro := trampoline(func() packet {
			return evalNode(lexicalEnv, f.Body)
		})

		if shouldEvalMacros {
			return bounce(func() packet {
				// This is executed in the environment of its application, not the
				// environment of its definition
				return evalNode(dynamicEnv, expandedMacro)
			})
		} else {
			return respond(expandedMacro)
		}
	} else {
		// Evaluate the body in the new lexical environment
		return bounce(func() packet {
			return evalNode(lexicalEnv, f.Body)
		})
	}
}

func checkSpecialArgs(name string, head ast.Node, args []ast.Node, paramCountMin int, paramCountMax int) {
	checkBuiltinArgs("Special form", name, head, args, paramCountMin, paramCountMax)
}

func checkPrimitiveArgs(name string, head ast.Node, args []ast.Node, paramCountMin int, paramCountMax int) {
	checkBuiltinArgs("Primitive", name, head, args, paramCountMin, paramCountMax)
}

func checkBuiltinArgs(builtinType string, name string, head ast.Node, args []ast.Node, paramCountMin int, paramCountMax int) {
	switch {
	case paramCountMax == -1:
		if !(paramCountMin <= len(args)) {
			panicEvalError(head, fmt.Sprintf(
				"%v '%v' expects at least %v argument(s), but was given %v",
				builtinType,
				name,
				paramCountMin,
				len(args)))
		}
	case paramCountMin == paramCountMax:
		if !(paramCountMin == len(args)) {
			panicEvalError(head, fmt.Sprintf(
				"%v '%v' expects %v argument(s), but was given %v",
				builtinType,
				name,
				paramCountMin,
				len(args)))
		}
	default:
		if !(paramCountMin <= len(args) && len(args) <= paramCountMax) {
			panicEvalError(head, fmt.Sprintf(
				"%v '%v' expects between %v and %v arguments, but was given %v",
				builtinType,
				name,
				paramCountMin,
				paramCountMax,
				len(args)))
		}
	}
}

func toSymbolName(n ast.Node) string {
	switch value := n.(type) {
	case *ast.Symbol:
		return value.Name
	}

	panic("Not a symbol: " + n.String())
}
