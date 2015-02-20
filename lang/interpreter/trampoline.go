package interpreter

import "vamos/lang/ast"

// Packet contains a thunk or a ast.Node.
// A packet is the result of the evaluation of a thunk.
type packet struct {
	Thunk thunk
	Node  ast.Node
}

// Bounce continues the trampolining session by placing a new thunk in the chain.
func bounce(t thunk) packet {
	return packet{Thunk: t}
}

// Respond exits a trampolining session by placing a ast.Node on the end of the
// chain.
func respond(n ast.Node) packet {
	return packet{Node: n}
}

type thunk func() packet

// Trampoline iteratively calls a chain of thunks until there is no next thunk,
// at which point it pulls the resulting ast.Node out of the packet and returns it.
func trampoline(currentThunk thunk) ast.Node {
	for currentThunk != nil {
		nextPacket := currentThunk()

		if nextPacket.Thunk != nil {
			currentThunk = nextPacket.Thunk
		} else {
			return nextPacket.Node
		}
	}

	return nil
}
