package parsing

import (
	"fmt"
	"vamos/lang/scanning"
)

type ParserError struct {
	Pos     scanning.TokenPosition
	Message string
}

// Implements the error interface
func (p *ParserError) Error() string {
	return fmt.Sprintf("Error at %v: %v", p.Pos, p.Message)
}

type ParserErrorList []*ParserError

func NewParserErrorList() ParserErrorList {
	return make(ParserErrorList, 0)
}

func (p *ParserErrorList) Add(pos scanning.TokenPosition, msg string) {
	*p = append(*p, &ParserError{pos, msg})
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
