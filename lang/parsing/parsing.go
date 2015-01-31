package parsing

import (
	"strconv"
	. "vamos/lang/ast"
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
	c := p.peek().Code
	if c == TcEOF || c == TcError {
		return true
	}

	return false
}

////////// Parsing

func parseNodes(p *parser, errors *ParserErrorList) []Node {
	var nodes []Node
	//TODO nodes := make([]Node, 0)
	for !p.inputEmpty() {
		nodes = append(nodes, parseAnnotatedNode(p, errors))
	}
	return nodes
}

func parseAnnotatedNode(p *parser, errors *ParserErrorList) AnnotatedNode {
	token := p.next()

	switch token.Code {
	case TcLeftParen:
		var list []Node
		for p.peek().Code != TcRightParen {
			if p.peek().Code == TcEOF {
				errors.Add(token.Pos, "Unbalanced parentheses")
				p.next()
				return &Symbol{Name: "nil"}
			}
			list = append(list, parseAnnotatedNode(p, errors))
		}
		p.next()
		return &List{Nodes: list}
	case TcRightParen:
		errors.Add(token.Pos, "Unbalanced parentheses")
	case TcNumber:
		return parseNumber(token, errors)
	case TcSymbol:
		return parseSymbol(token)
	case TcCaret:
		return parseAnnotation(p, errors)
	case TcSingleQuote:
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
	var list []Node
	list = append(list, &Symbol{Name: "quote"}, node)
	return &List{Nodes: list}
}

func parseNumber(t Token, errors *ParserErrorList) *Number {
	f, ferr := strconv.ParseFloat(t.Value, 64)

	if ferr != nil {
		errors.Add(t.Pos, "Invalid number: "+t.Value)
		return &Number{Value: 0.0}
	}

	return &Number{Value: f}
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
