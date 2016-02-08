update-element!
(def x '(1 2 3 4))
(update-element! x 0 10)
(def y '(1 (2 3) 4 5))
(def z (first (rest y)))
(update-element! (first (rest y)) 0 20)

(list x y z)
