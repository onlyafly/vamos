package lang

import (
	"fmt"
	"strconv"
	"unicode/utf8"
)

func Parse(input string, sourceName string) (Nodes, ParserErrorList) {
	s, _ := Scan(sourceName, input)
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
				errors.Add(token.Loc, "Unbalanced parentheses")
				p.next()
				return &NilNode{Location: token.Loc}
			}
			list = append(list, parseAnnotatedNode(p, errors))
		}
		p.next()
		return &ListNode{Nodes: list, Location: token.Loc}
	case TcRightParen:
		errors.Add(token.Loc, "Unbalanced parentheses")
	case TcNumber:
		return parseNumber(token, errors)
	case TcSymbol:
		return parseSymbol(token, errors)
	case TcString:
		return parseString(token, errors)
	case TcChar:
		return parseChar(token, errors)
	case TcCaret:
		return parseAnnotation(p, errors)
	case TcSingleQuote:
		return parseQuote(p, errors)
	default:
		errors.Add(token.Loc, "Unrecognized token: "+token.String())
	}

	return &NilNode{Location: token.Loc}
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
	return &ListNode{Nodes: list}
}

func parseNumber(t Token, errors *ParserErrorList) *Number {
	f, ferr := strconv.ParseFloat(t.Value, 64)

	if ferr != nil {
		errors.Add(t.Loc, "Invalid number: "+t.Value)
		return &Number{Value: 0.0, Location: t.Loc}
	}

	return &Number{Value: f, Location: t.Loc}
}

func parseSymbol(t Token, errors *ParserErrorList) AnnotatedNode {
	if t.Value == "nil" {
		return &NilNode{Location: t.Loc}
	}
	return &Symbol{Name: t.Value, Location: t.Loc}
}

func parseString(t Token, errors *ParserErrorList) *StringNode {
	content := t.Value[1 : len(t.Value)-1]
	return &StringNode{Value: content, Location: t.Loc}
}

func parseChar(t Token, errors *ParserErrorList) *CharNode {
	switch {
	case t.Value == "\\newline":
		return &CharNode{Value: '\n'}
	case len(t.Value) == 2:
		_, leadingSlashWidth := utf8.DecodeRuneInString(t.Value)
		r, _ := utf8.DecodeRuneInString(t.Value[leadingSlashWidth:])
		return &CharNode{Value: r}
	}

	errors.Add(t.Loc, fmt.Sprintf("Invalid character literal: %v", t.Value))
	return &CharNode{}
}

////////// Helper Functions

func ensureSymbol(n Node) *Symbol {
	if v, ok := n.(*Symbol); ok {
		return v
	}

	panic("Expected symbol: " + n.String())
}

func ensureList(n Node) *ListNode {
	if v, ok := n.(*ListNode); ok {
		return v
	}

	panic("Expected list: " + n.String())
}
