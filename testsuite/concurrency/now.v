(def len
  (proc (xs)
    (if (= xs '())
      0
      (+ 1 (len (rest xs))))))

(def n (now))

(len n)
