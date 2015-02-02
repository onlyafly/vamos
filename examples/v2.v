(defn lookup (exp env)
  (eval exp env))

(def else true)

(defn e2 (exp env)
  (if (atom? exp)
    (cond
      (symbol? exp) (lookup exp env)
      else          exp)
    'failed))
