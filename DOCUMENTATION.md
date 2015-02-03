# Vamos Language Documentation

* See http://hyperpolyglot.org/lisp for comparison of different Lisps.

## Literals

    "This is a string"

    42
    
    'this_is_a_quoted_symbol

    '()

## Special Forms

    ; This is a (single-line) comment

    (def x 4)

    (set! x 1)

    (quote x) --OR-- 'x

    (list 1 2 3)

    (if <BOOL> <THEN> <ELSE>)

    (cond
      <BOOL 1> <THEN 1>
      <BOOL 2> <THEN 2>
      <BOOL 3> <THEN 3>)

    (let (x 4
          y (+ 1 x))
      (* x y))

    (apply + '(1 3))

    (eval '(+ 1 2))

### Functions

    (fn (x y z)
      (+ x (+ y z)))

Variable number of arguments:

    (fn (x y &rest z)
      (+ x (+ y z)))

## Metaprogramming

    (def defn
      (macro (name args body)
        (list 'def name
          (list 'fn args
            body))))

    (macroexpand1 '(defn inc (a) (+ 1 a)))
    => (def inc (fn (a) (+ 1 a)))

    (function-params inc)
    => (n)

    (function-body inc)
    => (+ 1 n)

    (function-environment inc)
    => #environment<TopLevel>

## Built-in Functions

Math:

    +, -

Logical (on all types):

    =

Logical (on numbers):

    <, >

Lists:

    first, rest, list

Higher-order:

    apply

Other:

    (typeof 4)
    => number

### Evaluation and environments

    (current-environment)

## Boolean Values

False values: false (the symbol), which is also stored in false (the variable)
True values: true (the symbol), also stored in true (the variable)
