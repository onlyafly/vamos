(defproc start (n acc)
    (if (< n 0)
        (println acc)
        (start (- n 1) (+ 1 acc))))

(start 100000 0)