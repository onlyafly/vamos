(def fib
  (fn (n)
    (if (> n 1)
      (+ (fib (- n 1))
         (fib (- n 2)))
      n)))

(def fib-iter-help
  (fn (a b n)
    (if (> n 0)
      (fib-iter-help (+ a b)
                     a
                     (- n 1))
      b)))

(def fib-iter
  (fn (n)
    (fib-iter-help 1 0 n)))
