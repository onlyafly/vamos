package parsing

import (
	"strconv"
	"vamos/lang/ast"
	"vamos/lang/scanning"
)

func Parse(input string) *ast.Module {
	s, _ := scanning.Scan("parserCreated", input)
	p := &parser{s: s}

	return analyzeModule(parseModule(p))
}

////////// Parser

type parser struct {
	s              *scanning.Scanner
	lookahead      [2]scanning.Token // two-token lookahead
	lookaheadCount int
}

func (p *parser) next() scanning.Token {
	if p.lookaheadCount > 0 {
		p.lookaheadCount--
	} else {
		p.lookahead[0] = <-p.s.Tokens
	}
	return p.lookahead[p.lookaheadCount]
}

func (p *parser) backup() {
	p.lookaheadCount++
}

func (p *parser) peek() scanning.Token {
	if p.lookaheadCount > 0 {
		return p.lookahead[p.lookaheadCount-1]
	}

	p.lookaheadCount = 1
	p.lookahead[0] = <-p.s.Tokens
	return p.lookahead[0]
}

func (p *parser) inputEmpty() bool {
	if p.peek().Code == scanning.TC_EOF {
		return true
	}

	return false
}

////////// Parsing

func parseModule(p *parser) *ast.Module {
	nodes := make([]ast.Node, 0)
	for !p.inputEmpty() {
		nodes = append(nodes, parseNode(p))
	}
	return &ast.Module{nodes}
}

func parseNode(p *parser) ast.Node {
	token := p.next()

	switch {
	case token.Code == scanning.TC_LEFT_PAREN:
		list := make([]ast.Node, 0)
		for p.peek().Code != scanning.TC_RIGHT_PAREN {
			list = append(list, parseNode(p))
		}
		p.next()
		return &ast.List{list}
	case token.Code == scanning.TC_RIGHT_PAREN:
		panic("unexpected )")
	default:
		return parseAtom(token.Value)
	}

	return &ast.Symbol{"nil"}
}

func parseAtom(tokenValue string) ast.Expr {
	f, ferr := strconv.ParseFloat(tokenValue, 64)
	if ferr == nil {
		return &ast.Number{f}
	}

	return &ast.Symbol{tokenValue}
}

////////// Semantic Analysis

func analyzeModule(m *ast.Module) *ast.Module {
	nodes := make([]ast.Node, len(m.Nodes))

	for i, n := range m.Nodes {
		nodes[i] = analyzeTopLevelNode(n)
	}
	return &ast.Module{nodes}
}

func analyzeTopLevelNode(n ast.Node) ast.Decl {
	switch v := n.(type) {
	case *ast.List:
		return analyzeTopLevelList(v)
	default:
		panic("Unable to analyze declaration: " + n.String())
	}

	return nil
}

func analyzeTopLevelList(list *ast.List) ast.Decl {
	first := list.Nodes[0]

	switch v := first.(type) {
	case *ast.Symbol:
		switch v.Name {
		case "defn":
			return analyzeFunctionDecl(list.Nodes[1:])
		case "package":
			return analyzePackageDecl(list.Nodes[1:])
		}
	}

	panic("Unrecognized top-level list: " + list.String())
}

func analyzeFunctionDecl(nodes []ast.Node) *ast.FunctionDecl {
	functionNameSymbol := ensureSymbol(nodes[0])
	argumentsList := ensureList(nodes[1])
	body := nodes[2:]
	return &ast.FunctionDecl{functionNameSymbol, argumentsList, body}
}

func analyzePackageDecl(nodes []ast.Node) *ast.PackageDecl {
	nameSymbol := ensureSymbol(nodes[0])
	return &ast.PackageDecl{nameSymbol}
}

////////// Helper Functions

func ensureSymbol(n ast.Node) *ast.Symbol {
	if v, ok := n.(*ast.Symbol); ok {
		return v
	}

	panic("Expected symbol: " + n.String())
}

func ensureList(n ast.Node) *ast.List {
	if v, ok := n.(*ast.List); ok {
		return v
	}

	panic("Expected list: " + n.String())
}
