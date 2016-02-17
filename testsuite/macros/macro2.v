(def defproc
  (macro
    (proc (name args body)
      (list 'def name
        (list 'proc args
          body)))))

(defproc addem (a b) (+ a b))

(addem 100 1000)
