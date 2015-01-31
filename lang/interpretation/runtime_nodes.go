package interpretation

import . "vamos/lang/ast"

////////// Primitive

type primitiveFunction func(Env, []Node) Node

type Primitive struct {
	Name  string
	Value primitiveFunction
}

func (p *Primitive) String() string {
	return "#primitive<" + p.Name + ">"
}

func (p *Primitive) Children() []Node { return nil }
func (p *Primitive) isExpr() bool     { return true }
func (p *Primitive) Equals(n Node) bool {
	panicEvalError("Cannot compare the values of primitive procedures: " +
		p.String() + " and " + n.String())
	return false
}

////////// Function

type Function struct {
	Name       string
	Parameters []Node
	Body       Node
	ParentEnv  Env
}

func (f *Function) String() string {
	return "#function<" + f.Name + ">"
}

func (f *Function) Children() []Node { return nil }
func (f *Function) isExpr() bool     { return true }
func (f *Function) Equals(n Node) bool {
	panicEvalError("Cannot compare the values of functions: " +
		f.String() + " and " + n.String())
	return false
}
