package interpretation

import (
	"fmt"
	. "vamos/lang/ast"
)

// Eval evaluates a node in an environment.
func Eval(e Env, n Node) (result Node, err error) {
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
	case *List:
		return bounce(func() packet { return evalList(e, value, true) })
	default:
		panicEvalError("Unknown form to evaluate: " + value.String())
	}

	return respond(&Symbol{Name: "nil"})
}

func evalList(e Env, l *List, shouldEvalMacros bool) packet {
	elements := l.Nodes

	if len(elements) == 0 {
		panicEvalError("Empty list cannot be evaluated: " + l.String())
		return respond(nil)
	}

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
				return evalList(e, &List{Nodes: nodes}, true)
			}))
		case "def":
			return evalSpecialDef(e, args)
		case "eval":
			return evalSpecialEval(e, args)
		case "set!":
			name := toSymbolName(args[0])
			e.Update(name, trampoline(func() packet {
				return evalNode(e, args[1])
			}))
			return respond(&Symbol{Name: "nil"})
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
			panicEvalError("No matching cond clause: " + l.String())
		case "fn":
			return evalSpecialFn(e, args)
		case "macro":
			return evalSpecialMacro(e, args)
		case "macroexpand1":
			return evalSpecialMacroexpand1(e, args)
		case "quote":
			return respond(args[0])
		case "let":
			return bounce(func() packet {
				return evalLet(e, args)
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
			return evalFunctionApplication(value, arguments)
		})
	case *Macro:
		return bounce(func() packet {
			return evalMacroApplication(e, value, args, shouldEvalMacros)
		})
	default:
		panicEvalError("First item in list not a function: " + value.String())
	}

	return respond(&Symbol{Name: "nil"})
}

func evalFunctionApplication(f *Function, args []Node) packet {

	ensureArgsMatchParameters(f.Name, &args, &f.Parameters)

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

	// Save arguments into parameters
	for i, arg := range args {
		paramName := toSymbolName(f.Parameters[i])
		e.Set(paramName, arg)
	}

	return bounce(func() packet {
		return evalNode(e, f.Body)
	})
}

func evalMacroApplication(applicationEnv Env, m *Macro, args []Node, shouldEvalMacros bool) packet {
	macroResult := expandMacro(m, args)

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

func expandMacro(m *Macro, args []Node) Node {

	ensureArgsMatchParameters(m.Name, &args, &m.Parameters)

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

func evalLet(parentEnv Env, args []Node) packet {
	variableList := args[0]
	body := args[1]
	variableNodes := variableList.Children()

	e := NewMapEnv("let", parentEnv)

	// Evaluate variable assignments
	for i := 0; i < len(variableNodes); i += 2 {
		variable := variableNodes[i]
		expression := variableNodes[i+1]
		variableName := toSymbolName(variable)

		e.Set(variableName, trampoline(func() packet {
			return evalNode(e, expression)
		}))
	}

	// Evaluate body
	return bounce(func() packet {
		return evalNode(e, body)
	})
}

func ensureSpecialArgsCountEquals(formName string, args []Node, paramCount int) {
	if len(args) != paramCount {
		panicEvalError(fmt.Sprintf(
			"Special form '%v' expects %v argument(s), but was given %v",
			formName,
			paramCount,
			len(args)))
	}
}

func ensureSpecialArgsCountInRange(specialName string, args []Node, paramCountMin int, paramCountMax int) {
	if !(paramCountMin <= len(args) && len(args) <= paramCountMax) {
		panicEvalError(fmt.Sprintf(
			"Special form '%v' expects between %v and %v arguments, but was given %v",
			specialName,
			paramCountMin,
			paramCountMax,
			len(args)))
	}
}

func ensureArgsMatchParameters(procedureName string, args *[]Node, params *[]Node) {
	if len(*args) != len(*params) {
		panicEvalError(fmt.Sprintf(
			"Procedure '%v' expects %v argument(s), but was given %v",
			procedureName,
			len(*params),
			len(*args)))
	}
}

func panicEvalError(s string) {
	panic(NewEvalError(s))
}

func toSymbolName(n Node) string {
	switch value := n.(type) {
	case *Symbol:
		return value.Name
	}

	panic("Not a symbol: " + n.String())
}
