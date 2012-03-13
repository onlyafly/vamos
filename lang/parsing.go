package lang

import (
	"strconv"
)

func Parse(input string) Expression {
	return analyzeTopLevelExpression(parseExpression(tokenize(input)))
}

////////// Parsing

func parseExpression(tokens *TokenList) Expression {
	if tokens.empty() {
		panic("Unexpected EOF while parsing expression")
	}
	token := tokens.pop()

	switch {
	case "(" == token:
		list := make([]Expression, 0)
		for tokens.top() != ")" {
			list = append(list, parseExpression(tokens))
		}
		tokens.pop()
		return NewList(list)
	case ")" == token:
		panic("unexpected )")
	default:
		return parseAtom(token)
	}

	return NewSymbol("nil")
}

func parseAtom(token string) Expression {
	f, ferr := strconv.ParseFloat(token, 64)
	if ferr == nil {
		return NewNumber(f)
	}

	return NewSymbol(token)
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
