(def try
     (proc (guess x)
         (if (good-enough guess x)
             guess
           (try (improve guess x) x))))

(def sqrt (proc (x) (try 1 x)))

(def improve
     (proc (guess x)
         (average guess (/ x guess))))

(def good-enough
     (proc (guess x)
         (< (abs (- (square guess) x))
            0.001)))

(def abs
     (proc (x)
         (if (< x 0)
             (- 0 x)
           x)))

(def square
     (proc (x)
         (* x x)))

(def average
     (proc (x y)
         (/ (+ x y) 2)))

(sqrt 64)
