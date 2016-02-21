(def p (proc (a b c) (list 1 a 2 b 3 c 4)))
(apply p '(101 102 103))
