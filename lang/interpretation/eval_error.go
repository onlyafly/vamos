package interpretation

import (
	"fmt"
	. "vamos/lang/helpers"
)

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
