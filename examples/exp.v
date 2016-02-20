(def foo
  (let (x 1)
    (proc (y)
      (+ y x))))

(def foo +)

(defproc m ()
  (let (efoo (function-environment foo))
    (begin
      (println (foo 2))
      (eval '(update! x 2) efoo)
      (println (foo 2)))))
