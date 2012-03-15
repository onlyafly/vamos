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

////////// Token

type token struct {
	code  tokenCode
	value string
}

func (t token) String() string {
	switch t.code {
	case tkEOF:
		return "EOF"
	case tkError:
		return t.value
	}

	return fmt.Sprintf("%q", t.value)
}

////////// Token code

type tokenCode int

const (
	tkError tokenCode = iota
	tkLeftParen
	tkRightParen
	tkSymbol
	tkNumber
	tkEOF
)

const eof = -1

////////// Scanner struct

type scanner struct {
	name   string     // used only for error reports
	input  string     // the string being scanned
	start  int        // start position of this item
	pos    int        // current position in the input
	width  int        // width of last rune read from input
	tokens chan token // channel of scanned items
}

func (s *scanner) run() {
	for state := scanBegin; state != nil; {
		state = state(s)
	}
	close(s.tokens)
}

func (s *scanner) emit(code tokenCode) {
	s.tokens <- token{code, s.input[s.start:s.pos]}
	s.start = s.pos
}

func (s *scanner) next() (r rune) {
	if s.pos >= len(s.input) {
		s.width = 0
		return eof
	}
	r, s.width = utf8.DecodeRuneInString(s.input[s.pos:])
	s.pos += s.width
	return r
}

// ignore skips over the pending input before this point.
func (s *scanner) ignore() {
	s.start = s.pos
}

// backup steps back one rune.
// Can be called only once per call of next.
func (s *scanner) backup() {
	s.pos -= s.width
}

// peek returns but does not consume
// the next rune in the input.
func (s *scanner) peek() rune {
	r := s.next()
	s.backup()
	return r
}

// accept consumes the next rune
// if it's from the valid set.
func (s *scanner) accept(valid string) bool {
	if strings.IndexRune(valid, s.next()) >= 0 {
		return true
	}
	s.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (s *scanner) acceptRun(valid string) {
	for strings.IndexRune(valid, s.next()) >= 0 {
	}
	s.backup()
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (s *scanner) errorf(format string, args ...interface{}) stateFn {
	s.tokens <- token{
		tkError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

////////// Scanning

type stateFn func(*scanner) stateFn

func scan(name, input string) (*scanner, chan token) {
	s := &scanner{
		name:   name,
		input:  input,
		tokens: make(chan token),
	}
	go s.run()
	return s, s.tokens
}

func scanBegin(s *scanner) stateFn {
	for {
		switch r := s.next(); {
		case isSpace(r):
			s.ignore()
		case r == '(':
			s.emit(tkLeftParen)
		case r == ')':
			s.emit(tkRightParen)
		case r == '+' || r == '-' || '0' <= r && r <= '9':
			s.backup()
			return scanNumber
			/*
			 case isAlphaNumeric(r):
			 s.backup()
			 return scanSymbol
			*/
		case r == eof:
			break
		}
	}

	s.emit(tkEOF)
	return nil
}

func scanNumber(s *scanner) stateFn {
	// Optional leading sign.
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
	// Next thing mustn't be alphanumeric.
	if isAlphaNumeric(s.peek()) {
		s.next()
		return s.errorf("bad number syntax: %q",
			s.input[s.start:s.pos])
	}
	s.emit(tkNumber)
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
