(def bar
  (proc ()
    (println "bar!")))

(+ 9 (eval (read-string "(begin (bar) (+ 1 1))")))
