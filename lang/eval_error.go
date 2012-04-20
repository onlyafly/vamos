package lang

import (
	"fmt"
)

type EvalError string

// Implements the error interface
func (e EvalError) Error() string {
	return fmt.Sprintf("Evaluation error: %v", string(e))
}
