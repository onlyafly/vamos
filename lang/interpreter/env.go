package interpreter

import (
	"fmt"
	"github.com/onlyafly/vamos/lang/ast"
)

////////// Env

// Env represents an environment.
// An environment (AKA a scope) contains symbols that are in scope and which
// environment, if any, is the parent of this environment.
type Env interface {
	Set(name string, value ast.Node)
	Update(name string, value ast.Node) bool
	Get(name string) (ast.Node, bool)
	String() string
	Parent() Env
	Name() string
}

////////// MapEnv

// MapEnv is an implementation of an environment using a hash map.
type MapEnv struct {
	name    string
	symbols map[string]ast.Node
	parent  Env
}

// NewTopLevelMapEnv creates a new top-level envirxonment, which is initialized
// with the primitives.
func NewTopLevelMapEnv() *MapEnv {
	e := &MapEnv{
		name:    "TopLevel",
		symbols: make(map[string]ast.Node),
		parent:  nil,
	}

	initializePrimitives(e)

	return e
}

// NewMapEnv creates a new (non-top-level) environment.
func NewMapEnv(name string, parent Env) *MapEnv {
	return &MapEnv{
		name:    name,
		symbols: make(map[string]ast.Node),
		parent:  parent,
	}
}

// Set sets the initial value of a symbol.
func (e *MapEnv) Set(name string, value ast.Node) {
	if _, exists := e.symbols[name]; exists {
		panicEvalError(value, "Cannot set the initial value of a symbol again: "+name)
	} else {
		e.symbols[name] = value
	}
}

// Update updates the value of an existing symbol.
func (e *MapEnv) Update(name string, value ast.Node) bool {
	_, exists := e.symbols[name]

	if !exists {
		if e.Parent() == nil {
			return false
		} else {
			return e.Parent().Update(name, value)
		}
	}

	e.symbols[name] = value
	return true
}

// Get returns the value of a symbol.
func (e *MapEnv) Get(name string) (ast.Node, bool) {
	value, exists := e.symbols[name]

	if !exists {
		if e.Parent() == nil {
			return nil, false
		} else {
			return e.Parent().Get(name)
		}
	}

	return value, true
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
