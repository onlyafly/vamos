(load "vtest.v")

;;;;;;;;;; Tests

(deftest "Truthful values in 'if'"
  (= (list
       (if (quote true) 1 2)
       (if true 1 2))
     '(1 1)))

(runtests tests)
