# Vamos Language Documentation

* TODO: Move everything here to COMPARISON.html

### Routines: Functions, Procedures, Primitives, and Macros

- Procedures allow side-effects
- Procedures are pure functions, allowing no side effects, and only accepting immutable data structures
- Primitives are built-in procedures and functions
- Macros

A procedure:

    (proc (x y z)
      (+ x (+ y z)))

Variable number of arguments:

    (proc (x y &rest z)
      (+ x (+ y z)))

### Metaprogramming

    (def defproc
      (macro
        (proc (name args body)
          (list 'def name
            (list 'proc args
              body)))))

    (macroexpand1 '(defproc inc (a) (+ 1 a)))
    => (def inc (proc (a) (+ 1 a)))

    (routine-params inc)
    => (n)

    (routine-body inc)
    => (+ 1 n)

    (routine-environment inc)
    => #environment<TopLevel>

### Built-in Routines

Math:

    +, -

Logical (on all types):

    =

Logical (on numbers):

    <, >

Lists:

    list

ast.Collections (nil, lists, strings):

    cons, concat, first, rest, empty?

Higher-order:

    apply

Strings:

    (concat "abc" "de" "fgh")
    => "abcdefgh"

Other:

    (typeof 4)
    => number

### Concurrency

    (now)
    => (2015 10 01 20 45 16) ;; == 2015-10-01 8:45:16 PM

    ;; Sleep for 1 second
    (sleep 1000)
    => nil

    (go
      (sleep 1000)
      (println "Woah!"))
    (println "Cool!")
    => Cool!
    => Woah!

    (def c (chan))
    (go (send! c 42))
    (take! c)
    => 42

    (def c (chan))
    (go (send! c 42))
    (close! c)
    (take! c)
    (take! c) ; doesn't block, since the channel is closed
    => 42
    => nil

### Evaluation and environments

    (current-environment)

    (load "test.v")

    (eval '(+ 1 2))

    (eval '(+ x y) environment)

    (read-string "(+ 41 1)")
    => (+ 41 1)

## Boolean Values

False values: false (the symbol), which is also stored in false (the variable)
True values: true (the symbol), also stored in true (the variable)
