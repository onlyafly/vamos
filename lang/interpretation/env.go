package interpretation

import "fmt"
import . "vamos/lang/ast"

////////// Env

// Env represents an environment.
// An environment (AKA a scope) contains symbols that are in scope and which
// environment, if any, is the parent of this environment.
type Env interface {
	Set(name string, value Node)
	Update(name string, value Node)
	Get(name string) Node
	String() string
	Parent() Env
	Name() string
}

////////// MapEnv

// MapEnv is an implementation of an environment using a hash map.
type MapEnv struct {
	name    string
	symbols map[string]Node
	parent  Env
}

// NewTopLevelMapEnv creates a new top-level environment, which is initialized
// with the primitives.
func NewTopLevelMapEnv() *MapEnv {
	e := &MapEnv{
		name:    "TopLevel",
		symbols: make(map[string]Node),
		parent:  nil,
	}

	initializePrimitives(e)

	return e
}

// NewMapEnv creates a new (non-top-level) environment.
func NewMapEnv(name string, parent Env) *MapEnv {
	return &MapEnv{
		name:    name,
		symbols: make(map[string]Node),
		parent:  parent,
	}
}

// Set sets the initial value of a symbol.
func (e *MapEnv) Set(name string, value Node) {
	if _, exists := e.symbols[name]; exists {
		panicEvalError("Cannot redefine a name: " + name)
	} else {
		e.symbols[name] = value
	}
}

// Update updates the value of an existing symbol.
func (e *MapEnv) Update(name string, value Node) {
	if _, exists := e.symbols[name]; exists {
		e.symbols[name] = value
	} else {
		panicEvalError("Cannot update an undefined name: " + name)
	}
}

// Get returns the value of a symbol.
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

// Parent returns the parent environment.
func (e *MapEnv) Parent() Env {
	return e.parent
}

func (e *MapEnv) String() string {
	return fmt.Sprintf("%v:%v", e.name, e.symbols)
}

func (e *MapEnv) Name() string {
	return e.name
}
