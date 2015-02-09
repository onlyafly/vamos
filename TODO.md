# Plan for Vamos

## Vision

***Vamos is a fast Lisp that encourages a purely functional programming style an
has CSP-like concurrency semantics.***

Principles:

1. Influenced by Go's concurrency
2. Influenced by Clojure's take on a Lisp
   - http://hyperpolyglot.org/lisp
   - http://www.tryclj.com/
3. Well-documented code in a near literate style

## About the Vision

Quotes from the CSP paper:

* [Processes] may not communicate with each other by updating global variables.
* In parallel programming coroutines appear as a more fundamental program structure than subroutines, which can be regarded as a special case.
* [A coroutine] may use input commands to achieve the effect of "multiple entry points" ... [and be] used like a SIMULA class instance as a concrete representation for abstract data.

## Upcoming Goals

1. Build up a metacircular evaluator (v2.v)
2. A compiler written in Vamos, that translates Vamos code to Go code
   - Learn how to do it via Lisp in Small Pieces

## Todo

### Next

- make it possible for the test framework to test v2
  - add ability to load libs at runtime
  - factor test framework into its own lib
  -
- move if to prelude (see "if2")
- research lisp primitives

### Tech Debt

- How can the interpreter be reimplemented to allow call/cc
- Make use of the unused walk.go
- Clean up eval.go

### Language

- Associate source information with Function nodes (filename, line number, etc)
- Get ideas from comparison of different lisps at: http://hyperpolyglot.org/lisp
- Make use of unused annotation functionality (see test 0100)
- Make use of the Decl interface
- Improve macros by making them non-first-class (?)
- Dynamic binding

### Tooling

### Other
