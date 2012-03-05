package lang

import (
	"fmt"
)

////////// Environment

type Env struct {
	Entries map[string]Form
}

func NewEnv() *Env {
	entries := map[string]Form {
		"test": NewNumber(32.4),
		"+": NewNativeFunction(func(args []Form) Form {
			return NewNumber(toNumberValue(args[0]) + toNumberValue(args[1]))
		}),
	}
	return &Env{entries}
}

type NativeFunction struct {
	Value func([]Form) Form
}

func NewNativeFunction(fun func([]Form)Form) *NativeFunction {
	return &NativeFunction{fun}
}

func (self *NativeFunction) String() string {
	return "NATIVE"
}

func (self *NativeFunction) Type() string {
	return "nativeFunction"
}

////////// Evaluation

func Eval(form Form, env *Env) Form {
	switch value := form.(type) {
	case *Symbol:
		return env.Entries[value.Name]
	case *Number:
		return value
	case *List:
		return evalCall(value, env)
	default:
		panic("Unknown form to evaluate: " + value.String())
	}

	return NewSymbol("nil")
}

func evalCall(list *List, env *Env) Form {
	funName := toSymbolValue(list.Value[0])

	args := evalForms(list.Value[1:], env)
	
	switch value := env.Entries[funName].(type) {
	case *NativeFunction:
		return value.Value(args)
	default:
		panic("Unable to execute function: " + list.String())
	}

	return NewSymbol("nil")	
}

func evalForms(forms []Form, env *Env) []Form {
	fmt.Printf("evalForms: %v\n", forms)
	values := make([]Form, len(forms))
	for i, form := range forms {
		values[i] = Eval(form, env)
	}
	fmt.Printf("ending: %v\n", values)
	return values
}
