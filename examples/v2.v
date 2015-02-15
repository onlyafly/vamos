(load "vtest.v")

;; Inspired by Lisp in Small Pieces (LiSP)

(defn lookup (exp env)
  (eval exp env))

(defn get (l n)
  (if (= n 0)
    (first l)
    (get (rest l) (- n 1))))

(defn e2begin (exps env)
  (if (list? exps)
    (if (not (empty? (rest exps)))
      (begin
        (e2 (first exps) env)
        (e2begin (rest exps) env))
      (e2 (first exps) env))
    (list)))

(defn e2 (e env)
  (if (atom? e)

    ;; Handle atoms
    (cond
      (symbol? e) (lookup e env)
      (or (number? e) (string? e) (boolean? e)) e
      else 'CANNOT_EVALUATE)

    ;; Handle non-atoms
    (let (proc (first e))
      (cond
        (= proc 'quote)  (get e 1)
        (= proc 'if)     (if (e2 (get e 1) env)
                             (e2 (get e 2) env)
                             (e2 (get e 3) env))
        (= proc 'begin)  (e2begin (rest e) env)
        ;;TODO (= proc 'set!)   (update! (get e 1) env (e2 (get e 2) env))
        ;;TODO (= prod 'fn)     (make-function (get e 0) (get e 1) env)
        ;;TODO else (invoke (e2 (first e) env)
        ;;TODO              (evlis (rest e) env))
        else (list 'NOT_YET_IMPLEMENTED proc)
        ))))

(let (env (current-environment))
  (begin

    (defvtest "Atoms"
      (vt= (e2 '2 env)
           2)
      (vt= (e2 "test" env)
           "test"))

    (defvtest "Begin"
      (vt= (e2 '(begin 5 4) env)
           4))

  ))

(vt-start)
