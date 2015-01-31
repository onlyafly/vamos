(def defn
  (macro (name args body)
    (list 'def name
      (list 'fn args
        body))))

(defn addem (a b) (+ a b))

(addem 100 1000)
