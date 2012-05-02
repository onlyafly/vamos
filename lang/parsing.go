package lang

import (
	"strconv"
)

func Parse(input string) (Nodes, ParserErrorList) {
	s, _ := Scan("parserCreated", input)
	errorList := NewParserErrorList()
	p := &parser{s: s}

	nodes := parseNodes(p, &errorList)

	if errorList.Len() > 0 {
		return nil, errorList
	}

	return nodes, nil
}

////////// Parser

type parser struct {
	s              *Scanner
	lookahead      [2]Token // two-token lookahead
	lookaheadCount int
}

func (p *parser) next() Token {
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

func (p *parser) peek() Token {
	if p.lookaheadCount > 0 {
		return p.lookahead[p.lookaheadCount-1]
	}

	p.lookaheadCount = 1
	p.lookahead[0] = <-p.s.Tokens
	return p.lookahead[0]
}

func (p *parser) inputEmpty() bool {
	if p.peek().Code == TC_EOF {
		return true
	}

	return false
}

////////// Parsing

func parseNodes(p *parser, errors *ParserErrorList) []Node {
	nodes := make([]Node, 0)
	for !p.inputEmpty() {
		nodes = append(nodes, parseAnnotatedNode(p, errors))
	}
	return nodes
}

func parseAnnotatedNode(p *parser, errors *ParserErrorList) AnnotatedNode {
	token := p.next()

	switch token.Code {
	case TC_LEFT_PAREN:
		list := make([]Node, 0)
		for p.peek().Code != TC_RIGHT_PAREN {
			list = append(list, parseAnnotatedNode(p, errors))
		}
		p.next()
		return &List{Nodes: list}
	case TC_RIGHT_PAREN:
		errors.Add(token.Pos, "Unbalanced parentheses")
	case TC_NUMBER:
		return parseNumber(token, errors)
	case TC_SYMBOL:
		return parseSymbol(token)
	case TC_CARET:
		return parseAnnotation(p, errors)
	case TC_SINGLE_QUOTE:
		return parseQuote(p, errors)
	default:
		errors.Add(token.Pos, "Unrecognized token: "+token.String())
	}

	return &Symbol{Name: "nil"}
}

func parseAnnotation(p *parser, errors *ParserErrorList) AnnotatedNode {
	annotation := parseAnnotatedNode(p, errors)
	annotatee := parseAnnotatedNode(p, errors)
	annotatee.SetAnnotation(annotation)
	return annotatee
}

func parseQuote(p *parser, errors *ParserErrorList) AnnotatedNode {
	node := parseAnnotatedNode(p, errors)
	list := make([]Node, 0)
	list = append(list, &Symbol{Name: "quote"}, node)
	return &List{Nodes: list}
}

func parseNumber(t Token, errors *ParserErrorList) *Number {
	f, ferr := strconv.ParseFloat(t.Value, 64)
	if ferr == nil {
		return &Number{Value: f}
	} else {
		errors.Add(t.Pos, "Invalid number: "+t.Value)
	}

	return &Number{Value: 0.0}
}

func parseSymbol(t Token) *Symbol {
	return &Symbol{Name: t.Value}
}

////////// Helper Functions

func ensureSymbol(n Node) *Symbol {
	if v, ok := n.(*Symbol); ok {
		return v
	}

	panic("Expected symbol: " + n.String())
}

func ensureList(n Node) *List {
	if v, ok := n.(*List); ok {
		return v
	}

	panic("Expected list: " + n.String())
}
