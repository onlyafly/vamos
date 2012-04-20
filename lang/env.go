package lang

import ()

////////// Env

type Env interface {
	Set(name string, value Node)
	Get(name string) Node
}

////////// MapEnv

type MapEnv struct {
	symbols map[string]Node
}

func NewMapEnv() *MapEnv {
	e := &MapEnv{make(map[string]Node)}
	initializePrimitives(e) // TODO
	return e
}

func (e *MapEnv) Set(name string, value Node) {
	if _, exists := e.symbols[name]; exists {
		panicEvalError("Cannot redefine a name: " + name)
	} else {
		e.symbols[name] = value
	}
}

func (e *MapEnv) Get(name string) Node {
	value, exists := e.symbols[name]

	if !exists {
		panicEvalError("Name not defined: " + name)
	}

	return value
}
