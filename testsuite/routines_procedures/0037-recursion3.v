(def f
     (proc (x)
         (if (g x)
             (+ x 1)
           (+ x 10))))

(def g
     (proc (y)
         true))

(f 7)
