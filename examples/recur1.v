(let (f (fn (x) x)
      y (f 1)
      z (+ z 1))
  (begin
    (println y)
    (println z)))
