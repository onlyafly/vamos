package lang

import (
	"fmt"
)

// EvalError represents an error that occurs during evaluation.
type EvalError struct {
	Message string
}

func NewEvalError(message string) *EvalError {
	return &EvalError{message}
}

// Implements the error interface
func (e *EvalError) Error() string {
	return fmt.Sprintf("Evaluation error: %v", e.Message)
}
