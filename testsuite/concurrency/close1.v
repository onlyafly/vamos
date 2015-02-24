(def c (chan))
(go (send! c 42)
    (close! c))
(println (take! c) (take! c) (take! c))
