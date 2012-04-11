package parsing

import (
	"strconv"
	"vamos/lang/ast"
	"vamos/lang/scanning"
)

func Parse(input string) (*ast.Module, ParserErrorList) {
	s, _ := scanning.Scan("parserCreated", input)
	errorList := NewParserErrorList()
	p := &parser{s: s}

	parsedModule := parseModule(p, &errorList)
	analyzedModule := analyzeModule(parsedModule, &errorList)
	return analyzedModule, errorList
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

func parseModule(p *parser, errors *ParserErrorList) *ast.Module {
	nodes := make([]ast.Node, 0)
	for !p.inputEmpty() {
		nodes = append(nodes, parseBasicNode(p, errors))
	}
	return &ast.Module{nodes}
}

func parseBasicNode(p *parser, errors *ParserErrorList) ast.BasicNode {
	token := p.next()

	switch token.Code {
	case scanning.TC_LEFT_PAREN:
		list := make([]ast.Node, 0)
		for p.peek().Code != scanning.TC_RIGHT_PAREN {
			list = append(list, parseBasicNode(p, errors))
		}
		p.next()
		return &ast.List{Nodes: list}
	case scanning.TC_RIGHT_PAREN:
		errors.Add(token.Pos, "Unbalanced parentheses")
	case scanning.TC_NUMBER:
		return parseNumber(token, errors)
	case scanning.TC_SYMBOL:
		return parseSymbol(token)
	case scanning.TC_CARET:
		return parseAnnotation(p, errors)
	default:
		errors.Add(token.Pos, "Unrecognized token: "+token.String())
	}

	return &ast.Symbol{Name: "nil"}
}

func parseAnnotation(p *parser, errors *ParserErrorList) ast.BasicNode {
	annotation := parseBasicNode(p, errors)
	annotatee := parseBasicNode(p, errors)
	annotatee.SetAnnotation(annotation)
	return annotatee
}

func parseNumber(t scanning.Token, errors *ParserErrorList) *ast.Number {
	f, ferr := strconv.ParseFloat(t.Value, 64)
	if ferr == nil {
		return &ast.Number{Value: f}
	} else {
		errors.Add(t.Pos, "Invalid number: "+t.Value)
	}

	return &ast.Number{Value: 0.0}
}

func parseSymbol(t scanning.Token) *ast.Symbol {
	return &ast.Symbol{Name: t.Value}
}

////////// Semantic Analysis

func analyzeModule(m *ast.Module, errors *ParserErrorList) *ast.Module {
	nodes := make([]ast.Node, len(m.Nodes))

	for i, n := range m.Nodes {
		nodes[i] = analyzeTopLevelDecl(n)
	}
	return &ast.Module{nodes}
}

func analyzeTopLevelDecl(n ast.Node) ast.Decl {
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
	return &ast.FunctionDecl{
		Name:      functionNameSymbol,
		Arguments: argumentsList,
		Body:      body,
	}
}

func analyzePackageDecl(nodes []ast.Node) *ast.PackageDecl {
	nameSymbol := ensureSymbol(nodes[0])
	return &ast.PackageDecl{Name: nameSymbol}
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
