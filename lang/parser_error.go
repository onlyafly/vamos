package lang

import (
	"fmt"
	"vamos/lang/token"
)

////////// ParserError

type ParserError struct {
	Loc     *token.Location
	Message string
}

// Implements the error interface
func (pe *ParserError) Error() string {
	return fmt.Sprintf("Error (line %v): %v", pe.Loc.Line, pe.Message)
}

////////// ParserErrorList

// ParserErrorList is a list of ParserError pointers.
// Implements the error interface.
type ParserErrorList []*ParserError

func NewParserErrorList() ParserErrorList {
	return make(ParserErrorList, 0)
}

func (p *ParserErrorList) Add(loc *token.Location, msg string) {
	*p = append(*p, &ParserError{loc, msg})
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
