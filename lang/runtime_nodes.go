package lang

////////// EnvNode

type EnvNode struct {
	Env Env
}

func NewEnvNode(e Env) *EnvNode {
	return &EnvNode{Env: e}
}

func (en *EnvNode) Name() string {
	return en.Env.Name()
}

func (en *EnvNode) String() string {
	return "#environment<" + en.Env.Name() + ">"
}

func (en *EnvNode) Children() []Node    { return nil }
func (en *EnvNode) isExpr() bool        { return true }
func (en *EnvNode) TypeName() string    { return "environment" }
func (en *EnvNode) Loc() *TokenLocation { return nil }
func (en *EnvNode) Equals(n Node) bool {
	panicEvalError(n, "Cannot compare the values of environments: "+
		en.String()+" and "+n.String())
	return false
}

////////// Primitive

type primitiveFunction func(Env, []Node) Node

type Primitive struct {
	Name  string
	Value primitiveFunction
}

func (p *Primitive) String() string {
	return "#primitive<" + p.Name + ">"
}

func (p *Primitive) Children() []Node    { return nil }
func (p *Primitive) isExpr() bool        { return true }
func (p *Primitive) TypeName() string    { return "primitive" }
func (p *Primitive) Loc() *TokenLocation { return nil }
func (p *Primitive) Equals(n Node) bool {
	panicEvalError(n, "Cannot compare the values of primitive procedures: "+
		p.String()+" and "+n.String())
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

func (f *Function) Children() []Node    { return nil }
func (f *Function) isExpr() bool        { return true }
func (f *Function) TypeName() string    { return "function" }
func (f *Function) Loc() *TokenLocation { return nil }
func (f *Function) Equals(n Node) bool {
	panicEvalError(n, "Cannot compare the values of functions: "+
		f.String()+" and "+n.String())
	return false
}

////////// Macro

type Macro struct {
	Name       string
	Parameters []Node
	Body       Node
	ParentEnv  Env
}

func (m *Macro) String() string {
	return "#macro<" + m.Name + ">"
}

func (m *Macro) Children() []Node    { return nil }
func (m *Macro) isExpr() bool        { return true }
func (m *Macro) TypeName() string    { return "macro" }
func (m *Macro) Loc() *TokenLocation { return nil }
func (m *Macro) Equals(n Node) bool {
	panicEvalError(n, "Cannot compare the values of macros: "+
		m.String()+" and "+n.String())
	return false
}
