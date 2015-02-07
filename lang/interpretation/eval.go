package interpretation

import (
	"fmt"
	"io"
	. "vamos/lang/ast"
	. "vamos/lang/helpers"
)

var writer io.Writer

// Eval evaluates a node in an environment.
func Eval(e Env, n Node, w io.Writer) (result Node, err error) {
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

	startThunk := func() packet {
		return evalNode(e, n)
	}

	return trampoline(startThunk), nil
}

func evalEachNode(e Env, ns []Node) []Node {
	result := make([]Node, len(ns))
	for i, n := range ns {
		evalNodeThunk := func() packet {
			return evalNode(e, n)
		}
		result[i] = trampoline(evalNodeThunk)
	}
	return result
}

func evalNode(e Env, n Node) packet {

	switch value := n.(type) {
	case *Number:
		return respond(value)
	case *Symbol:
		return respond(e.Get(value.Name))
	case *StringNode:
		return respond(value)
	case *CharNode:
		return respond(value)
	case *ListNode:
		return bounce(func() packet { return evalList(e, value, true) })
	case *NilNode:
		return respond(&NilNode{})
	default:
		panicEvalError(n, "Unknown form to evaluate: "+value.String())
	}

	return respond(&NilNode{})
}

func evalList(e Env, l *ListNode, shouldEvalMacros bool) packet {
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
	case *Symbol:
		switch value.Name {
		case "apply":
			f := args[0]
			l := toListValue(trampoline(func() packet {
				return evalNode(e, args[1])
			}))
			nodes := append([]Node{f}, l.Nodes...)
			return respond(trampoline(func() packet {
				return evalList(e, &ListNode{Nodes: nodes}, true)
			}))
		case "def":
			return evalSpecialDef(e, head, args)
		case "eval":
			return evalSpecialEval(e, head, args)
		case "update!":
			name := toSymbolName(args[0])
			e.Update(name, trampoline(func() packet {
				return evalNode(e, args[1])
			}))
			return respond(&NilNode{})
		case "if":
			predicate := toBooleanValue(trampoline(func() packet {
				return evalNode(e, args[0])
			}))
			if predicate {
				return bounce(func() packet {
					return evalNode(e, args[1])
				})
			}
			return bounce(func() packet {
				return evalNode(e, args[2])
			})
		case "cond":
			for i := 0; i < len(args); i += 2 {
				predicate := toBooleanValue(trampoline(func() packet {
					return evalNode(e, args[i])
				}))

				if predicate {
					return bounce(func() packet {
						return evalNode(e, args[i+1])
					})
				}
			}
			panicEvalError(head, "No matching cond clause: "+l.String())
		case "fn":
			return evalSpecialFn(e, head, args)
		case "macro":
			return evalSpecialMacro(e, head, args)
		case "macroexpand1":
			return evalSpecialMacroexpand1(e, head, args)
		case "quote":
			return respond(args[0])
		case "let":
			return bounce(func() packet {
				return evalSpecialLet(e, head, args)
			})
		case "begin":
			return bounce(func() packet {
				return evalSpecialBegin(e, head, args)
			})
		}
	}

	headNode := trampoline(func() packet {
		return evalNode(e, head)
	})

	switch value := headNode.(type) {
	case *Primitive:
		f := value.Value
		return respond(f(e, evalEachNode(e, args)))
	case *Function:
		arguments := evalEachNode(e, args)

		return bounce(func() packet {
			return evalFunctionApplication(value, head, arguments)
		})
	case *Macro:
		return bounce(func() packet {
			return evalMacroApplication(e, value, head, args, shouldEvalMacros)
		})
	default:
		panicEvalError(head, "First item in list not a function: "+value.String())
	}

	return respond(&NilNode{})
}

func evalFunctionApplication(f *Function, head Node, args []Node) packet {

	isVariableNumberOfParams := false
	for _, param := range f.Parameters {
		paramName := toSymbolName(param)
		if paramName == "&rest" {
			isVariableNumberOfParams = true
		}
	}
	if !isVariableNumberOfParams {
		ensureArgsMatchParameters(f.Name, head, &args, &f.Parameters)
	}

	e := NewMapEnv(f.Name, f.ParentEnv)

	// TODO
	/*
		print(
			"evalFunctionApplication:\n   name=",
			e.String(), "\n   body=",
			f.Body.String(), "\n   parent=",
			f.ParentEnv.String(), "\n   args=",
			fmt.Sprintf("%v", args), "\n   isTail=",
			fmt.Sprintf("%v", isTail), "\n")
	*/

	// Map arguments to parameters
	isMappingRestArgs := false
	iarg := 0
	for iparam, param := range f.Parameters {
		paramName := toSymbolName(param)
		if isMappingRestArgs {
			restArgs := args[iarg:]
			restList := NewList(restArgs)
			e.Set(paramName, restList)
		} else if paramName == "&rest" {
			isMappingRestArgs = true
		} else {
			e.Set(paramName, args[iparam])
			iarg++
		}
	}

	return bounce(func() packet {
		return evalNode(e, f.Body)
	})
}

func evalMacroApplication(applicationEnv Env, m *Macro, head Node, args []Node, shouldEvalMacros bool) packet {
	macroResult := expandMacro(m, head, args)

	if shouldEvalMacros {
		return bounce(func() packet {
			// This is executed in the environment of its application, not the
			// environment of its definition
			return evalNode(applicationEnv, macroResult)
		})
	} else {
		return respond(macroResult)
	}
}

func expandMacro(m *Macro, head Node, args []Node) Node {
	ensureArgsMatchParameters(m.Name, head, &args, &m.Parameters)

	e := NewMapEnv(m.Name, m.ParentEnv)

	// TODO

	/*
		print(
			"evalMacroApplication:\n   name=",
			e.String(), "\n   body=",
			m.Body.String(), "\n   parent=",
			m.ParentEnv.String(), "\n   args=",
			fmt.Sprintf("%v", args), "\n")
	*/

	// Save arguments into parameters
	for i, arg := range args {
		paramName := toSymbolName(m.Parameters[i])
		e.Set(paramName, arg)
	}

	macroResult := trampoline(func() packet {
		return evalNode(e, m.Body)
	})

	return macroResult
}

func ensureSpecialArgsCountEquals(formName string, head Node, args []Node, paramCount int) {
	if len(args) != paramCount {
		panicEvalError(head, fmt.Sprintf(
			"Special form '%v' expects %v argument(s), but was given %v",
			formName,
			paramCount,
			len(args)))
	}
}

func ensureSpecialArgsCountInRange(specialName string, head Node, args []Node, paramCountMin int, paramCountMax int) {
	if !(paramCountMin <= len(args) && len(args) <= paramCountMax) {
		panicEvalError(head, fmt.Sprintf(
			"Special form '%v' expects between %v and %v arguments, but was given %v",
			specialName,
			paramCountMin,
			paramCountMax,
			len(args)))
	}
}

func ensureArgsMatchParameters(procedureName string, head Node, args *[]Node, params *[]Node) {
	if len(*args) != len(*params) {
		panicEvalError(head, fmt.Sprintf(
			"Procedure '%v' expects %v argument(s), but was given %v",
			procedureName,
			len(*params),
			len(*args)))
	}
}

func panicEvalError(n Node, s string) {
	var loc *TokenLocation
	if n != nil {
		loc = n.Loc()
	}
	panic(NewEvalError(s, loc))
}

func toSymbolName(n Node) string {
	switch value := n.(type) {
	case *Symbol:
		return value.Name
	}

	panic("Not a symbol: " + n.String())
}
