package interpretation

import (
	. "vamos/lang/ast"
)

func toListValue(n Node) *List {
	switch value := n.(type) {
	case *List:
		return value
	}

	panicEvalError("Expression is not a list: " + n.String())
	return nil
}

func toNumberValue(n Node) float64 {
	switch value := n.(type) {
	case *Number:
		return value.Value
	}

	panicEvalError("Expression is not a number: " + n.String())
	return 0.0
}

func toSymbolValue(exp Expr) string {
	switch value := exp.(type) {
	case *Symbol:
		return value.Name
	}

	panicEvalError("Expression is not a symbol: " + exp.String())
	return ""
}

func toBooleanValue(n Node) bool {
	switch value := n.(type) {
	case *Symbol:
		switch value.Name {
		case "true":
			return true
		case "false":
			return false
		}
	}

	// All other values are treated as false
	return false
}
