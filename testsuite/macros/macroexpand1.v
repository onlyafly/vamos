(def defn
  (macro (name args body)
    (list 'def name
      (list 'fn args
        body))))

(macroexpand1 '(defn addem (a b) (+ a b)))
