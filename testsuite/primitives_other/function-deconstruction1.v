(def inc
  (let (add (proc (a b)
              (+ a b)))
    (proc (n)
      (add n 1))))

(def rebuild-inc
  (macro
    (proc (name)
      (list 'def name
        (list 'proc (function-params inc)
          (function-body inc))))))

(eval '(rebuild-inc inc2) (function-environment inc))

(list
  (inc 3)
  (function-params inc)
  (function-body inc)
  (function-environment inc)
  (eval '(inc2 5) (function-environment inc)))
