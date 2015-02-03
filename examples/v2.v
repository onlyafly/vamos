(defn lookup (exp env)
  (eval exp env))

(defn e2 (exp env)
  (if (atom? exp)
    (cond
      (symbol? exp)      (lookup exp env)
      (or (number? exp)
          (boolean? exp)) exp
      else               'CANNOT_EVALUATE)
    'NOT_YET_IMPLEMENTED))
