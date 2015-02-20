package interpreter

import (
	"fmt"
	"vamos/lang/ast"
	"vamos/lang/token"
)

// EvalError represents an error that occurs during evaluation.
type EvalError struct {
	Message  string
	location *token.Location
}

func NewEvalError(message string, location *token.Location) *EvalError {
	return &EvalError{message, location}
}

// Implements the error interface
func (e *EvalError) Error() string {
	if e.location != nil {
		return fmt.Sprintf("Evaluation error (%v: %v): %v", e.location.Filename, e.location.Line, e.Message)
	}

	return fmt.Sprintf("Evaluation error: %v", e.Message)
}

func panicEvalError(n ast.Node, s string) {
	var loc *token.Location
	if n != nil {
		loc = n.Loc()
	}
	panic(NewEvalError(s, loc))
}
