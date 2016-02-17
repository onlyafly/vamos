(def foo (proc ()
  (current-environment)))

(def bar (proc ()
  (current-environment)))

(def efoo (foo))
(def ebar (bar))

(eval '(def a 100) efoo)
(eval '(def a 1000) ebar)

(list
  (eval '(+ a 1) efoo)
  (eval '(+ a 1) ebar))
