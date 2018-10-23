(defproc start (n acc)
    ;(println n)
    (if (< n 0)
        (__stacktrace)
        (+ (start (- n 1) (+ 1 acc)) 0)))

(start 100000 0)