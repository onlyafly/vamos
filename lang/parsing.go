package lang

import (
	"strconv"
)

func Parse(input string) Expression {
	s, _ := Scan("parserCreated", input)
	p := &parser{s: s}

	return analyzeTopLevelExpression(parseExpression(p))
}

////////// Parser

type parser struct {
	s *scanner
	lookahead [2]token // two-token lookahead
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
		return p.lookahead[p.lookaheadCount - 1]
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

func parseExpression(p *parser) Expression {
	if p.inputEmpty() {
		panic("Unexpected EOF while parsing expression")
	}
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

func analyzeTopLevelExpression(e Expression) Expression {
	switch v := e.(type) {
	case *List:
		return analyzeTopLevelList(v)
	default:
		panic("Unrecognized top-level expression: " + e.String())
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
