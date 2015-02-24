package parser

import (
	"fmt"
	"strconv"
	"unicode/utf8"
	"vamos/lang/ast"
)

// Parse accepts a string and the name of the source of the code, and returns
// the Vamos nodes therein, along with a list of any errors found.
func Parse(input string, sourceName string) (ast.Nodes, ParserErrorList) {
	s, _ := Scan(sourceName, input)
	errorList := NewParserErrorList()
	p := &parser{s: s}

	s.errorHandler = func(t Token, message string) {
		errorList.Add(t.Loc, message)
	}

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

func parseNodes(p *parser, errors *ParserErrorList) []ast.Node {
	var nodes []ast.Node
	for !p.inputEmpty() {
		nodes = append(nodes, parseAnnotatedNode(p, errors))
	}
	return nodes
}

func parseAnnotatedNode(p *parser, errors *ParserErrorList) ast.AnnotatedNode {
	token := p.next()

	switch token.Code {
	case TcError:
		errors.Add(token.Loc, "Error token: "+token.String())
	case TcLeftParen:
		var list []ast.Node
		for p.peek().Code != TcRightParen {
			if p.peek().Code == TcEOF {
				errors.Add(token.Loc, "Unbalanced parentheses")
				p.next()
				return &ast.Nil{Location: token.Loc}
			}
			list = append(list, parseAnnotatedNode(p, errors))
		}
		p.next()
		return &ast.List{Nodes: list, Location: token.Loc}
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

	return &ast.Nil{Location: token.Loc}
}

func parseAnnotation(p *parser, errors *ParserErrorList) ast.AnnotatedNode {
	annotation := parseAnnotatedNode(p, errors)
	annotatee := parseAnnotatedNode(p, errors)
	annotatee.SetAnnotation(annotation)
	return annotatee
}

func parseQuote(p *parser, errors *ParserErrorList) ast.AnnotatedNode {
	node := parseAnnotatedNode(p, errors)
	var list []ast.Node
	list = append(list, &ast.Symbol{Name: "quote"}, node)
	return &ast.List{Nodes: list}
}

func parseNumber(t Token, errors *ParserErrorList) *ast.Number {
	f, ferr := strconv.ParseFloat(t.Value, 64)

	if ferr != nil {
		errors.Add(t.Loc, "Invalid number: "+t.Value)
		return &ast.Number{Value: 0.0, Location: t.Loc}
	}

	return &ast.Number{Value: f, Location: t.Loc}
}

func parseSymbol(t Token, errors *ParserErrorList) ast.AnnotatedNode {
	if t.Value == "nil" {
		return &ast.Nil{Location: t.Loc}
	}
	return &ast.Symbol{Name: t.Value, Location: t.Loc}
}

func parseString(t Token, errors *ParserErrorList) *ast.Str {
	content := t.Value[1 : len(t.Value)-1]
	return &ast.Str{Value: content, Location: t.Loc}
}

func parseChar(t Token, errors *ParserErrorList) *ast.Char {
	switch {
	case t.Value == "\\newline":
		return &ast.Char{Value: '\n'}
	case len(t.Value) == 2:
		_, leadingSlashWidth := utf8.DecodeRuneInString(t.Value)
		r, _ := utf8.DecodeRuneInString(t.Value[leadingSlashWidth:])
		return &ast.Char{Value: r}
	}

	errors.Add(t.Loc, fmt.Sprintf("Invalid character literal: %v", t.Value))
	return &ast.Char{}
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
