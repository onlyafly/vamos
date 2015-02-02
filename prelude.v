(def defn
  (macro (name args body)
    (list 'def name
      (list 'fn args
        body))))

(defn not (b)
  (cond
    (= b false) true
    (= b true)  false
    true        false))

(defn list? (n)
  (= (typeof n) 'list))

(defn symbol? (n)
  (= (typeof n) 'symbol))

(defn atom? (n)
  (not (list? n)))

'(loaded prelude version 2)
