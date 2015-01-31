(def defn
  (macro (name args body)
    (list 'def name
      (list 'fn args
        body))))
