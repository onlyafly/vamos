(def bar
  (fn ()
    (println "bar!")))

(+ 9 (eval (read-string "(begin (bar) (+ 1 1))")))
