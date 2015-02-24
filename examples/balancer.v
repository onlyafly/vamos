(def num-workers 42)

(defn sendLotsOfWork (c)
  (defn loop (i)
    (if (> i 0)
      (begin
        (send! c (list i (* i 2) 0))
        (loop (- i 1)))
      nil))
  (loop 100)
  (send! c 'DONE)
  )

(defn receiveLotsOfResults (c)
  (go (defn loop ()
        (let (w (take! c))
          (if (= w 'DONE)
            nil
            (begin
              (println "Received:" w)
              (loop)))))
      (loop))

  (sleep 100000))

(defn worker (cin cout)
  (let (w (take! cin))
    (if (= w 'DONE)
      nil
      (let (x (first w)
            y (first (rest w))
            z (* x y))
        (begin
          (println "Sleeping for: " z)
          (sleep z)
          (send! cout (list x y z))
          (worker cin cout))))))

(defn main ()
  (let (cin (chan)
        cout (chan))
    (begin

      (defn loop (i)
        (if (> i 0)
          (begin
            (go (worker cin cout))
            (loop (- i 1)))
          nil))
      (loop num-workers)

      (go (sendLotsOfWork cin))
      (receiveLotsOfResults cout))))

(main)
