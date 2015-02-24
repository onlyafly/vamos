(def num-workers 50)
(def work-count 100)

(defn sendLotsOfWork (cin)
  (defn loop (i)
    (if (< i work-count)
      (begin
        (send! cin (list i (* i 2) 0))
        (loop (+ i 1)))
      nil))
  (loop 0)
  (close! cin)
  )

(defn receiveLotsOfResults (cout)
  (go (defn loop ()
        (let (w (take! cout))
          (if w
            (begin
              (println "Received:" w)
              (loop))
            nil)))
      (loop))

  (println "Sleeping for 10 seconds...")
  (sleep 10000))

(defn worker (cin cout)
  (let (w (take! cin))
    (if w
      (let (x (first w)
            y (first (rest w))
            z (* x y))
        (begin
          (sleep z)
          (send! cout (list x y z))
          (worker cin cout)))
      nil)))

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
