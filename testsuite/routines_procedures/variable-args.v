(def foo
  (proc (x &rest xs)
    (list 'x= x 'xs= xs)))

(foo 1 2 3 4)
