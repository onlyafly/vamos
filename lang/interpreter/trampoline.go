package interpreter

import "vamos/lang/ast"

// A packet represents the continuation of a sequence of computations.
// It contains either a Next or a Node, but not both.
// If it contains a Next, the thunk is the next computation to execute.
// If it contains a Node, the trampolining session is over and the Node represents the result.
type packet struct {
	Next   thunk
	Result ast.Node
}

// Bounce continues the trampolining session by placing a new thunk in the chain.
func bounce(t thunk) packet {
	return packet{Next: t}
}

// Respond exits a trampolining session by placing a ast.Node on the end of the
// chain.
func respond(n ast.Node) packet {
	return packet{Result: n}
}

type thunk func() packet

// Trampoline iteratively calls a chain of thunks until there is no next thunk,
// at which point it pulls the resulting ast.Node out of the packet and returns it.
func trampoline(currentThunk thunk) ast.Node {
	for currentThunk != nil {
		nextPacket := currentThunk()

		if nextPacket.Next != nil {
			currentThunk = nextPacket.Next
		} else {
			return nextPacket.Result
		}
	}

	return nil
}
