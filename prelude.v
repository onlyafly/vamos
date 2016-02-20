;;;;;;;;;; Procedures

(def defproc
  (macro
    (proc (name args &rest exps)
      (list 'def name
        (list 'proc args
          (cons 'begin exps))))))

(def defmacro
  (macro
    (proc (name args body)
      (list 'def name
        (list 'macro
          (list 'proc args
            body))))))

;;;;;;;;;; Math

(defproc <= (a b)
  (or (< a b) (= a b)))

(defproc >= (a b)
  (or (> a b) (= a b)))

;;;;;;;;;; Logic

(def else true)

(defproc binary-or (a b)
  (cond
    (= a true) true
    (= b true) true
    else       false))

(defproc binary-and (a b)
  (if (= a true)
    (if (= b true)
      true
      false)
    false))

(defproc or (&rest xs)
  (foldl binary-or false xs))

(defproc and (&rest xs)
  (foldl binary-and true xs))

(defproc not (b)
  (cond
    (= b false) true
    (= b true)  false
    else        false))

;;;;;;;;;; Higher Order Procedures

(defproc foldl (f init xs)
  (if (= xs '())
    init
    (foldl f
           (f init (first xs))
           (rest xs))))

(defproc reverse (xs)
  (foldl (proc (acc x) (cons x acc)) '() xs))

(defproc map (f l)
  (let (loop (proc (accum xs)
               (if (empty? xs)
                 accum
                 (loop (cons (f (first xs)) accum)
                       (rest xs)))))
    (loop '() (reverse l))))

;;;;;;;;;; Type Predicates

(defproc list? (n)
  (= (typeof n) 'list))

(defproc char? (n)
  (= (typeof n) 'char))

(defproc symbol? (n)
  (= (typeof n) 'symbol))

(defproc number? (n)
  (= (typeof n) 'number))

(defproc procedure? (n)
  (= (typeof n) 'procedure))

(defproc macro? (n)
  (= (typeof n) 'macro))

(defproc environment? (n)
  (= (typeof n) 'environment))

(defproc primitive? (n)
  (= (typeof n) 'primitive))

(defproc string? (n)
  (= (typeof n) 'string))

(defproc atom? (n)
  (not (list? n)))

(defproc empty? (n)
  (cond (= n '()) true
        (= n "") true
        (= n nil) true
        else false))

(defproc boolean? (n)
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

;; TODO Naive implementation
(defproc len (xs)
  (if (empty? xs)
    0
    (+ 1 (len (rest xs)))))

;;;;;;;;;;

"Prelude version 2016-02-12"
