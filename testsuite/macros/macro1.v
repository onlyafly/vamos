(def firstm
  (macro
    (fn (a b)
      (list (first a)
            (first (rest b))
            100))))
            
(firstm (+ 1 2) (+ 3 4))
