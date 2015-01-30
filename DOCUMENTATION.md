# Vamos Language Documentation

* See http://hyperpolyglot.org/lisp for comparison of different Lisps.

## Special Forms

(def x 4)

(set! x 1)

(fn (x y z) (+ x (+ y z)))

(quote x) --OR-- 'x

(list 1 2 3)

(if <BOOL> <THEN> <ELSE>)

(cond
    (isTheAnswer 42
     (findResult?) 'foo))

(let (x 4
      y (+ 1 x))
  (* x y))

(apply + '(1 3))

## Built-in Functions

Math: +, -

Logical (on all types): =

Logical (on numbers): <, >

Lists: first, rest, list

Higher-order: apply

## Boolean Values

False values: false (the symbol), which is also stored in false (the variable)
True values: true (the symbol), also stored in true (the variable)
