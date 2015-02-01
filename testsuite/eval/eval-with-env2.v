(def capture (fn (a b)
  (current-environment)))

(def e (capture 1 2))

(eval '(+ a b) e)
