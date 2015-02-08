package lang

func toListValue(n Node) *ListNode {
	switch value := n.(type) {
	case *ListNode:
		return value
	}

	panicEvalError(n, "Expression is not a list: "+n.String())
	return nil
}

func toNumberValue(n Node) float64 {
	switch value := n.(type) {
	case *Number:
		return value.Value
	}

	panicEvalError(n, "Expression is not a number: "+n.String())
	return 0.0
}

func toSymbolValue(n Node) string {
	switch value := n.(type) {
	case *Symbol:
		return value.Name
	}

	panicEvalError(n, "Expression is not a symbol: "+n.String())
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
