package lang

import "fmt"

// EvalError represents an error that occurs during evaluation.
type EvalError struct {
	Message  string
	location *TokenLocation
}

func NewEvalError(message string, location *TokenLocation) *EvalError {
	return &EvalError{message, location}
}

// Implements the error interface
func (e *EvalError) Error() string {
	if e.location != nil {
		return fmt.Sprintf("Evaluation error (line %v): %v", e.location.Line, e.Message)
	}

	return fmt.Sprintf("Evaluation error: %v", e.Message)
}

func panicEvalError(n Node, s string) {
	var loc *TokenLocation
	if n != nil {
		loc = n.Loc()
	}
	panic(NewEvalError(s, loc))
}
