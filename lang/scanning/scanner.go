/*
See "Lexical Scanning in Go" by Rob Pike for the basic theory behind this
module: http://www.youtube.com/watch?v=HxaD_trXwRE
*/

package scanning

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

////////// Token

type Token struct {
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
	TC_EOF
)

const eof = -1

////////// Scanner struct

type Scanner struct {
	name   string     // used only for error reports
	input  string     // the string being scanned
	start  int        // start position of this item
	pos    int        // current position in the input
	width  int        // width of last rune read from input
	Tokens chan Token // channel of scanned items
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
	s.Tokens <- Token{code, s.input[s.start:s.pos]}
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

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (s *Scanner) errorf(format string, args ...interface{}) stateFn {
	s.Tokens <- Token{
		TC_ERROR,
		fmt.Sprintf(format, args...),
	}
	return nil
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
		case r == '+' || r == '-' || '0' <= r && r <= '9':
			s.backup()
			return scanNumber
			/*
			 case isAlphaNumeric(r):
			 s.backup()
			 return scanSymbol
			*/
		case isSymbolic(r):
			s.backup()
			return scanSymbol
		case r == eof:
			break Outer
		default:
			fmt.Printf("scanBegin default: %#v\n", r)
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
		return s.errorf("bad number syntax: %q",
			s.input[s.start:s.pos])
	}
	s.emit(TC_NUMBER)

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
	case r == '?':
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
