package lang

import (
	"fmt"
)

type ParserError struct {
	Pos     TokenPosition
	Message string
}

// Implements the error interface
func (p *ParserError) Error() string {
	return fmt.Sprintf("Error at %v: %v", p.Pos, p.Message)
}

// Implements the error interface
type ParserErrorList []*ParserError

func NewParserErrorList() ParserErrorList {
	return make(ParserErrorList, 0)
}

func (p *ParserErrorList) Add(pos TokenPosition, msg string) {
	*p = append(*p, &ParserError{pos, msg})
}

func (p ParserErrorList) Error() string {
	return p.String()
}

func (p ParserErrorList) Len() int {
	return len(p)
}

func (p ParserErrorList) String() (s string) {
	for i, e := range p {
		s += e.Error()

		if i != len(p)-1 {
			s += "\n"
		}
	}

	return s
}
