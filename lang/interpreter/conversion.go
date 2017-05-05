package interpreter

import "github.com/onlyafly/vamos/lang/ast"

func toListValue(n ast.Node) *ast.List {
	switch value := n.(type) {
	case *ast.List:
		return value
	}

	panicEvalError(n, "Expression is not a list: "+n.String())
	return nil
}

func toNumberValue(n ast.Node) float64 {
	switch value := n.(type) {
	case *ast.Number:
		return value.Value
	}

	panicEvalError(n, "Expression is not a number: "+n.String())
	return 0.0
}

func toSymbolValue(n ast.Node) string {
	switch value := n.(type) {
	case *ast.Symbol:
		return value.Name
	}

	panicEvalError(n, "Expression is not a symbol: "+n.String())
	return ""
}

func toBooleanValue(n ast.Node) bool {
	switch value := n.(type) {
	case *ast.Symbol:
		switch value.Name {
		case "false":
			return false
		}
	case *ast.Nil:
		return false
	}

	// All other values are treated as true
	return true
}
