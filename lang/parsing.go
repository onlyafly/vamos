package lang

import (
	"strconv"
)

func Parse(input string) *Module {
	s, _ := Scan("parserCreated", input)
	p := &parser{s: s}

	return analyzeModule(parseModule(p))
}

////////// Parser

type parser struct {
	s              *scanner
	lookahead      [2]token // two-token lookahead
	lookaheadCount int
}

func (p *parser) next() token {
	if p.lookaheadCount > 0 {
		p.lookaheadCount--
	} else {
		p.lookahead[0] = <-p.s.tokens
	}
	return p.lookahead[p.lookaheadCount]
}

func (p *parser) backup() {
	p.lookaheadCount++
}

func (p *parser) peek() token {
	if p.lookaheadCount > 0 {
		return p.lookahead[p.lookaheadCount-1]
	}

	p.lookaheadCount = 1
	p.lookahead[0] = <-p.s.tokens
	return p.lookahead[0]
}

func (p *parser) inputEmpty() bool {
	if p.peek().code == tkEOF {
		return true
	}

	return false
}

////////// Parsing

func parseModule(p *parser) *Module {
	es := make([]Expression, 0)
	for !p.inputEmpty() {
		es = append(es, parseExpression(p))
	}
	return &Module{es}
}

func parseExpression(p *parser) Expression {
	token := p.next()

	switch {
	case token.code == tkLeftParen:
		list := make([]Expression, 0)
		for p.peek().code != tkRightParen {
			list = append(list, parseExpression(p))
		}
		p.next()
		return NewList(list)
	case token.code == tkRightParen:
		panic("unexpected )")
	default:
		return parseAtom(token.value)
	}

	return NewSymbol("nil")
}

func parseAtom(tokenValue string) Expression {
	f, ferr := strconv.ParseFloat(tokenValue, 64)
	if ferr == nil {
		return NewNumber(f)
	}

	return NewSymbol(tokenValue)
}

////////// Semantic Analysis

func analyzeModule(m *Module) *Module {
	topLevelExpressions := make([]Expression, len(m.Expressions))

	for i, e := range m.Expressions {
		topLevelExpressions[i] = analyzeTopLevelExpression(e)
	}
	return &Module{topLevelExpressions}
}

func analyzeTopLevelExpression(e Expression) Expression {
	switch v := e.(type) {
	case *List:
		return analyzeTopLevelList(v)
	default:
		panic("Unable to analyze definition: " + e.String())
	}

	return nil
}

func analyzeTopLevelList(list *List) Expression {
	firstExpression := list.Value[0]

	switch v := firstExpression.(type) {
	case *Symbol:
		switch v.Name {
		case "defn":
			return analyzeFunctionDefinition(list.Value[1:])
		case "package":
			return analyzePackageDefinition(list.Value[1:])
		}
	}

	panic("Unrecognized top-level list: " + list.String())
}

func analyzeFunctionDefinition(es []Expression) *FunctionDefinition {
	functionNameSymbol := ensureSymbol(es[0])
	argumentsList := ensureList(es[1])
	body := es[2:]
	return NewFunctionDefinition(functionNameSymbol, argumentsList, body)
}

func analyzePackageDefinition(es []Expression) *PackageDefinition {
	nameSymbol := ensureSymbol(es[0])
	return NewPackageDefinition(nameSymbol)
}

////////// Helper Functions

func ensureSymbol(e Expression) *Symbol {
	if v, ok := e.(*Symbol); ok {
		return v
	}

	panic("Expected symbol: " + e.String())
}

func ensureList(e Expression) *List {
	if v, ok := e.(*List); ok {
		return v
	}

	panic("Expected list: " + e.String())
}
