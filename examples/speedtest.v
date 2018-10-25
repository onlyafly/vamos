(defproc start (n acc)
    ;(println n)
    (if (< n 0)
        acc
        (let (result (start (- n 1) (+ 1 acc)))
            result)))

;(start 100000 0)