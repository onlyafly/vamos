(def c (chan))
(go (send! c 42))
(take! c)
