(def inc
  (let (add (proc (a b)
              (+ a b)))
    (proc (n)
      (add n 1))))

(def rebuild-inc
  (macro
    (proc (name)
      (list 'def name
        (list 'proc (routine-params inc)
          (routine-body inc))))))

(eval '(rebuild-inc inc2) (routine-environment inc))

(list
  (inc 3)
  (routine-params inc)
  (routine-body inc)
  (routine-environment inc)
  (eval '(inc2 5) (routine-environment inc)))
