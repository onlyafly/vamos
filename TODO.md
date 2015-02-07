# Todo

- collection prims on strings
- write a test suite in Vamos for Vamos, to enable testing the prelude
- move if to prelude (see "if2")
- research lisp primitives

## Goal

- A compiler written in Vamos, that translates Vamos code to Go code
  - Learn how to do it via Lisp in Small Pieces

## Tech Debt

- How can the interpreter be reimplemented to allow call/cc
- Make use of the unused walk.go
- Clean up eval.go

## Language

- Associate source information with Function nodes (filename, line number, etc)
- Get ideas from comparison of different lisps at: http://hyperpolyglot.org/lisp
- Make use of unused annotation functionality (see test 0100)
- Make use of the Decl interface

## Tool

## Other

- Build up a metacircular evaluator (v2)
