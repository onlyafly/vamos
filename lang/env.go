package lang

import ()

type Env interface {
	Set(name string, value Node)
	Get(name string) Node
}

type MapEnv struct {
	symbols map[string]Node
}

func NewMapEnv() *MapEnv {
	return &MapEnv{make(map[string]Node)}
}

func (e *MapEnv) Set(name string, value Node) {
	if _, exists := e.symbols[name]; exists {
		panic("Cannot redefine a name: " + name)
	} else {
		e.symbols[name] = value
	}
}

func (e *MapEnv) Get(name string) Node {
	value, exists := e.symbols[name]

	if !exists {
		panic("Name not defined: " + name)
	}

	return value
}
