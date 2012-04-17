/*
See "Lexical Scanning in Go" by Rob Pike for the basic theory behind this
module: http://www.youtube.com/watch?v=HxaD_trXwRE
*/

package lang

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

////////// TokenPosition

type TokenPosition int

////////// Token

type Token struct {
	Pos   TokenPosition
	Code  TokenCode
	Value string
}

func (t Token) String() string {
	switch t.Code {
	case TC_EOF:
		return "EOF"
	case TC_ERROR:
		return t.Value
	}

	return fmt.Sprintf("%v", t.Value)
}

////////// Token code

type TokenCode int

const (
	TC_ERROR TokenCode = iota
	TC_LEFT_PAREN
	TC_RIGHT_PAREN
	TC_SYMBOL
	TC_NUMBER
	TC_CARET
	TC_EOF
)

const eof = -1

type ErrorHandler func(position int, message string)

////////// Scanner struct

type Scanner struct {
	name   string     // used only for error reports
	input  string     // the string being scanned
	start  int        // start position of this item
	pos    int        // current position in the input
	width  int        // width of last rune read from input
	Tokens chan Token // channel of scanned items

	// Error handling
	errorCount   int
	errorHandler ErrorHandler
}

func (s *Scanner) String() string {
	return fmt.Sprintf("<scanner remaining=%#v>", s.input[s.start:s.pos])
}

func (s *Scanner) run() {
	for state := scanBegin; state != nil; {
		state = state(s)
	}
	close(s.Tokens)
}

func (s *Scanner) emit(code TokenCode) {
	s.Tokens <- Token{
		Pos:   TokenPosition(s.start),
		Code:  code,
		Value: s.input[s.start:s.pos],
	}
	s.start = s.pos
}

func (s *Scanner) next() (r rune) {
	if s.pos >= len(s.input) {
		s.width = 0
		r = eof
		return
	}
	r, s.width = utf8.DecodeRuneInString(s.input[s.pos:])
	s.pos += s.width
	return
}

// ignore skips over the pending input before this point.
func (s *Scanner) ignore() {
	s.start = s.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (s *Scanner) backup() {
	s.pos -= s.width
}

// peek returns but does not consume
// the next rune in the input.
func (s *Scanner) peek() rune {
	r := s.next()
	s.backup()
	return r
}

// accept consumes the next rune
// if it's from the valid set.
func (s *Scanner) accept(valid string) bool {
	if strings.IndexRune(valid, s.next()) >= 0 {
		return true
	}
	s.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (s *Scanner) acceptRun(valid string) {
	for strings.IndexRune(valid, s.next()) >= 0 {
	}
	s.backup()
}

func (s *Scanner) emitErrorf(format string, args ...interface{}) {
	if s.errorHandler != nil {
		message := fmt.Sprintf(format, args...)
		s.errorHandler(s.start, message)
	}

	s.Tokens <- Token{
		Pos:   TokenPosition(s.start),
		Code:  TC_ERROR,
		Value: s.input[s.start:s.pos],
	}
	s.start = s.pos

	s.errorCount++
}

////////// Scanning

type stateFn func(*Scanner) stateFn

func Scan(name, input string) (*Scanner, chan Token) {
	s := &Scanner{
		name:   name,
		input:  input,
		Tokens: make(chan Token),
	}
	go s.run()
	return s, s.Tokens
}

func scanBegin(s *Scanner) stateFn {
Outer:
	for {
		switch r := s.next(); {
		case isSpace(r):
			s.ignore()
		case r == '(':
			s.emit(TC_LEFT_PAREN)
		case r == ')':
			s.emit(TC_RIGHT_PAREN)
		case '0' <= r && r <= '9':
			s.backup()
			return scanNumber
		case r == '+' || r == '-':
			rnext := s.next()
			if '0' <= rnext && rnext <= '9' {
				s.backup()
				s.backup()
				return scanNumber
			} else {
				s.backup()
				return scanSymbol
			}
		case r == '^':
			s.emit(TC_CARET)
		case isSymbolic(r):
			s.backup()
			return scanSymbol
		case r == eof:
			break Outer
		default:
			s.emitErrorf("unrecognized character")
		}
	}

	s.emit(TC_EOF)
	return nil
}

func scanSymbol(s *Scanner) stateFn {
	for isSymbolic(s.next()) {
	}
	s.backup()
	s.emit(TC_SYMBOL)
	return scanBegin
}

func scanNumber(s *Scanner) stateFn {
	// Optional leading sign
	s.accept("+-")

	// Is it hex?
	digits := "0123456789"
	if s.accept("0") && s.accept("xX") {
		digits = "0123456789abcdefABCDEF"
	}
	s.acceptRun(digits)
	if s.accept(".") {
		s.acceptRun(digits)
	}
	if s.accept("eE") {
		s.accept("+-")
		s.acceptRun("0123456789")
	}

	// Is it imaginary?
	s.accept("i")

	// Next thing must not be alphanumeric
	if isAlphaNumeric(s.peek()) {
		s.next()
		s.emitErrorf("bad number syntax: %q", s.input[s.start:s.pos])
	} else {
		s.emit(TC_NUMBER)
	}

	return scanBegin
}

////////// Helpers

func isAlphaNumeric(r rune) bool {
	switch {
	case '0' <= r && r <= '9':
		return true
	case 'a' <= r && r <= 'z':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	}

	return false
}

func isSymbolic(r rune) bool {
	switch {
	case '0' <= r && r <= '9':
		return true
	case 'a' <= r && r <= 'z':
		return true
	case 'A' <= r && r <= 'Z':
		return true
	case r == '?' || r == '+' || r == '-':
		return true
	}

	return false
}

func isSpace(r rune) bool {
	switch r {
	case ' ':
		return true
	case '\t':
		return true
	case '\r':
		return true
	case '\n':
		return true
	}

	return false
}
