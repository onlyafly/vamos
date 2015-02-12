;;;;;;;;;; Procedures

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

;;;;;;;;;; Logic

(def else true)

(defn binary-or (a b)
  (cond
    (= a true) true
    (= b true) true
    else       false))

(defn binary-and (a b)
  (if (= a true)
    (if (= b true)
      true
      false)
    false))

(defn or (&rest xs)
  (fold binary-or false xs))

(defn and (&rest xs)
  (fold binary-and true xs))

(defn not (b)
  (cond
    (= b false) true
    (= b true)  false
    else        false))

;;;;;;;;;; Higher Order Functions

(defn fold (f init xs)
  (if (= xs '())
    init
    (fold f
          (f init (first xs))
          (rest xs))))

;;;;;;;;;; Type Predicates

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

(defn string? (n)
  (= (typeof n) 'string))

(defn atom? (n)
  (not (list? n)))

(defn empty? (n)
  (cond (= n '()) true
        (= n "") true
        (= n nil) true
        else false))

(defn boolean? (n)
  (cond (= n true) true
        (= n false) true
        else false))

;; (if (= a b) (typeof a) (typeof b))
;; =>
;; (cond (= a b) (typeof a)
;;       true    (typeof b))
(defmacro if2 (condition consequent alternative)
  (list 'cond condition consequent
              true      alternative))

;;;;;;;;;;

"Prelude version 2015-02-08"
