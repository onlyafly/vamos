(def add
     (fn (x y)
         (if (= x 0)
             y
           (add (- x 1) (+ y 1)))))

(add 5 37)
