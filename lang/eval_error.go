package lang

import (
	"fmt"
)

// EvalError represents an error that occurs during evaluation.
type EvalError string

// Implements the error interface
func (e EvalError) Error() string {
	return fmt.Sprintf("Evaluation error: %v", string(e))
}
