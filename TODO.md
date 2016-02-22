# Plan for Vamos

## Vision

***Vamos is a fast Lisp that encourages a purely functional programming style and has CSP-like concurrency semantics.***

Principles:

1. Influenced by Go's concurrency
   - https://play.golang.org/
2. Influenced by Clojure's take on a Lisp
   - http://hyperpolyglot.org/lisp
   - http://clojure.org/lisps (see chart on page)
   - http://www.tryclj.com/
3. Well-documented code in a near literate style

## About the Vision

Quotes from the CSP paper:

* [Processes] may not communicate with each other by updating global variables.
* In parallel programming coroutines appear as a more fundamental program structure than subroutines, which can be regarded as a special case.
* [A coroutine] may use input commands to achieve the effect of "multiple entry points" ... [and be] used like a SIMULA class instance as a concrete representation for abstract data.

## Todo

### Next

- Allow escapes of double quotes in strings
- Make ast.Str.String() escape double quotes
- move if to prelude (see "prelude::if2") (?)
- Complete concurrency functionality

### Tech Debt

- How can the interpreter be reimplemented to allow call/cc, see Queinnec's LiSP §3
- Make use of the unused walk.go

### Future Ideas

- Add module system
- Associate source information with Procedure nodes (filename, line number, etc)
- Get ideas from comparison of different lisps at: http://hyperpolyglot.org/lisp
- Make use of unused annotation functionality (see test 0100)
- Improve macros by making them non-first-class (?)
- Types and type inference (?)
- Reader macros
- Dynamic scoping, see Queinnec's LiSP §2.5.1
- Exceptions to replace panic functionality
