(def inc
  (let (add (fn (a b)
              (+ a b)))
    (fn (n)
      (add n 1))))

(def rebuild-inc
  (macro
    (fn (name)
      (list 'def name
        (list 'fn (function-params inc)
          (function-body inc))))))

(eval '(rebuild-inc inc2) (function-environment inc))

(list
  (inc 3)
  (function-params inc)
  (function-body inc)
  (function-environment inc)
  (eval '(inc2 5) (function-environment inc)))
