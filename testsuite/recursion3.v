(def f
     (fn (x)
         (if (g x)
             (+ x 1)
           (+ x 10))))

(def g
     (fn (y)
         true))

(f 7)
