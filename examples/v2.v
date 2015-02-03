(defn lookup (exp env)
  (eval exp env))

(defn e2 (exp env)
  (if (atom? exp)
    ;; Handle atoms
    (cond
      (symbol? exp) (lookup exp env)
      (or (number? exp) (string? exp) (boolean? exp)) exp
      else 'CANNOT_EVALUATE)
    ;; Handle non-atoms
    (let (proc (first exp))
      (cond
        (= proc 'quote) (first (rest exp))
        (= proc 'if)    (if (e2)

        ))
    'NOT_YET_IMPLEMENTED))
