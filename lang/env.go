package lang

import ()

////////// Env

type Env interface {
	Set(name string, value Node)
	Get(name string) Node
	Parent() Env
}

////////// MapEnv

type MapEnv struct {
	symbols map[string]Node
	parent  Env
}

func NewTopLevelMapEnv() *MapEnv {
	e := &MapEnv{
		symbols: make(map[string]Node),
		parent:  nil,
	}

	initializePrimitives(e)

	return e
}

func NewMapEnv(parent Env) *MapEnv {
	return &MapEnv{
		symbols: make(map[string]Node),
		parent:  parent,
	}
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
		if e.Parent() == nil {
			panicEvalError("Name not defined: " + name)
		} else {
			return e.Parent().Get(name)
		}
	}

	return value
}

func (e *MapEnv) Parent() Env {
	return e.parent
}
