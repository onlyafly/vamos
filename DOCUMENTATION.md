# Vamos Language Documentation

* See http://hyperpolyglot.org/lisp for comparison of different Lisps.

### Literals

Nil:

    nil

Character:

    \a
    \newline

String:

    "This is a string"

Number:

    42

Symbol:

    'this_is_a_quoted_symbol
    true => 'true
    false => 'false

List:

    '()

### Special Forms

    ; This is a (single-line) comment

    (def x 4)

    (update! x 1)

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

    (begin
      (update! x 10)
      (number? x))
    => true

### Functions

    (fn (x y z)
      (+ x (+ y z)))

Variable number of arguments:

    (fn (x y &rest z)
      (+ x (+ y z)))

### Metaprogramming

    (def defn
      (macro
        (fn (name args body)
          (list 'def name
            (list 'fn args
              body)))))

    (macroexpand1 '(defn inc (a) (+ 1 a)))
    => (def inc (fn (a) (+ 1 a)))

    (function-params inc)
    => (n)

    (function-body inc)
    => (+ 1 n)

    (function-environment inc)
    => #environment<TopLevel>

### Built-in Functions

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

IO:

    (println "Test")

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
