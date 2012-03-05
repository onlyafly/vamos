package lang

import (
	"strconv"
)

////////// Parsing

func Parse(input string) Form {
	return readFrom(tokenize(input))
}

func readFrom(tokens *TokenList) Form {
	if tokens.empty() {
		panic("Unexpected EOF while reading")
	}
	token := tokens.pop()

	switch {
	case "(" == token:
		list := make([]Form, 0)
		for tokens.top() != ")" {
			list = append(list, readFrom(tokens))
		}
		tokens.pop()
		return NewList(list)
	case ")" == token:
		panic("unexpected )")
	default:
		return atom(token)
	}

	return NewSymbol("nil")
}

func atom(token string) Form {
	f, ferr := strconv.ParseFloat(token, 64)
	if ferr == nil {
		return NewNumber(f)
	}

	return NewSymbol(token)
}
