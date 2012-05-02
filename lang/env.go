package lang

import (
	"fmt"
)

////////// Env

type Env interface {
	Set(name string, value Node)
	Update(name string, value Node)
	Get(name string) Node
	String() string
	Parent() Env
}

////////// MapEnv

type MapEnv struct {
	name    string
	symbols map[string]Node
	parent  Env
}

func NewTopLevelMapEnv() *MapEnv {
	e := &MapEnv{
		name:    "TopLevel",
		symbols: make(map[string]Node),
		parent:  nil,
	}

	initializePrimitives(e)

	return e
}

func NewMapEnv(name string, parent Env) *MapEnv {
	return &MapEnv{
		name:    name,
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

func (e *MapEnv) Update(name string, value Node) {
	if _, exists := e.symbols[name]; exists {
		e.symbols[name] = value
	} else {
		panicEvalError("Cannot update an undefined name: " + name)
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

func (e *MapEnv) String() string {
	return fmt.Sprintf("%v:%v", e.name, e.symbols)
}
