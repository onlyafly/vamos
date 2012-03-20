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
	es := make([]ast.Expression, 0)
	for !p.inputEmpty() {
		es = append(es, parseExpression(p))
	}
	return &ast.Module{es}
}

func parseExpression(p *parser) ast.Expression {
	token := p.next()

	switch {
	case token.Code == scanning.TC_LEFT_PAREN:
		list := make([]ast.Expression, 0)
		for p.peek().Code != scanning.TC_RIGHT_PAREN {
			list = append(list, parseExpression(p))
		}
		p.next()
		return ast.NewList(list)
	case token.Code == scanning.TC_RIGHT_PAREN:
		panic("unexpected )")
	default:
		return parseAtom(token.Value)
	}

	return ast.NewSymbol("nil")
}

func parseAtom(tokenValue string) ast.Expression {
	f, ferr := strconv.ParseFloat(tokenValue, 64)
	if ferr == nil {
		return ast.NewNumber(f)
	}

	return ast.NewSymbol(tokenValue)
}

////////// Semantic Analysis

func analyzeModule(m *ast.Module) *ast.Module {
	topLevelExpressions := make([]ast.Expression, len(m.Expressions))

	for i, e := range m.Expressions {
		topLevelExpressions[i] = analyzeTopLevelExpression(e)
	}
	return &ast.Module{topLevelExpressions}
}

func analyzeTopLevelExpression(e ast.Expression) ast.Expression {
	switch v := e.(type) {
	case *ast.List:
		return analyzeTopLevelList(v)
	default:
		panic("Unable to analyze definition: " + e.String())
	}

	return nil
}

func analyzeTopLevelList(list *ast.List) ast.Expression {
	firstExpression := list.Value[0]

	switch v := firstExpression.(type) {
	case *ast.Symbol:
		switch v.Name {
		case "defn":
			return analyzeFunctionDefinition(list.Value[1:])
		case "package":
			return analyzePackageDefinition(list.Value[1:])
		}
	}

	panic("Unrecognized top-level list: " + list.String())
}

func analyzeFunctionDefinition(es []ast.Expression) *ast.FunctionDefinition {
	functionNameSymbol := ensureSymbol(es[0])
	argumentsList := ensureList(es[1])
	body := es[2:]
	return ast.NewFunctionDefinition(functionNameSymbol, argumentsList, body)
}

func analyzePackageDefinition(es []ast.Expression) *ast.PackageDefinition {
	nameSymbol := ensureSymbol(es[0])
	return ast.NewPackageDefinition(nameSymbol)
}

////////// Helper Functions

func ensureSymbol(e ast.Expression) *ast.Symbol {
	if v, ok := e.(*ast.Symbol); ok {
		return v
	}

	panic("Expected symbol: " + e.String())
}

func ensureList(e ast.Expression) *ast.List {
	if v, ok := e.(*ast.List); ok {
		return v
	}

	panic("Expected list: " + e.String())
}
