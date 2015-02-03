(def defn
  (macro (name args body)
    (list 'def name
      (list 'fn args
        body))))

(def defmacro
  (macro (name args body)
    (list 'def name
      (list 'macro args
        body))))

(def else true)

(defn simple-or (a b)
  (cond
    (= a true) true
    (= b true) true
    else       false))

(defn simple-and (a b)
  (if (= a true)
    (if (= b true)
      true
      false)
    false))

(defn not (b)
  (cond
    (= b false) true
    (= b true)  false
    else        false))

(defn fold (f init xs)
  (if (= xs '())
    init
    (fold f
          (f init (first xs))
          (rest xs))))

(defn list? (n)
  (= (typeof n) 'list))

(defn symbol? (n)
  (= (typeof n) 'symbol))

(defn number? (n)
  (= (typeof n) 'number))

(defn function? (n)
  (= (typeof n) 'function))

(defn macro? (n)
  (= (typeof n) 'macro))

(defn environment? (n)
  (= (typeof n) 'environment))

(defn primitive? (n)
  (= (typeof n) 'primitive))

(defn atom? (n)
  (not (list? n)))

'(loaded prelude version 2)
