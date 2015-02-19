(def get (fn (l n)
  (if (= n 0)
    (first l)
    (get (rest l) (- n 1)))))

(def start (now))

(sleep 2000)

(def end (now))

(if (= (get start 4) (get end 4))
  ;; Minute has not elapsed
  (- (get end 5)
     (get start 5))
  ;; Minute has elapsed, good enough
  2)
