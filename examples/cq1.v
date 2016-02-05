;; cq1
;; From Chapter 1 of "Lisp in Small Pieces" by Christian Queinnec
;; Updated 2016-02-05

(load "vtest.v")

(defn wrong (msg exp)
  ;; TODO should be equivalent to a panic?
  (panic "WRONG:" msg ":" exp))

(defn get (l n)
  (if (= n 0)
    (first l)
    (get (rest l) (- n 1))))

(defn proc? (p)
  (cond (function? p)  true
        (primitive? p) true
        else           false))

(defn qeval (e env)
  (if (atom? e)

    ;; Handle atoms
    (cond
      (symbol? e) (lookup e env)
      (or (number? e) (string? e) (boolean? e) (char? e)) e
      else (wrong "cannot evaluate" e))

    ;; Handle non-atoms
    (let (proc (first e))
      (cond
        (= proc 'quote)     (get e 1)
        (= proc 'if)        (if (qeval (get e 1) env)
                              (qeval (get e 2) env)
                              (qeval (get e 3) env))
        (= proc 'begin)     (qbegin (rest e) env)
        (= proc 'update!)   (qupdate! (get e 1) env (qeval (get e 2) env))
        (= proc 'fn)        (make-function (get e 0) (get e 1) env)
        else                (qapply (qeval (first e) env)
                                    (evlis (rest e) env))
        ))))

;; Environment is the closest thing to an Alist that Vamos supports:
;; ((a 1) (b 2) (c 3))
(defn lookup (id env)
  (if (empty? env)
    (wrong "no such binding" id)
    (if (= (get (first env) 0) id)
      (get (first env) 1)
      (lookup id (rest env)))))

;; Extend an environment env with a list of variables var and values val
(defn extend (env vars vals)
  (cond
    (empty? vars) (if (empty? vals)
                    env
                    (wrong "too few variables" (list vars vals)))
    (empty? vals) (wrong "too many variables" (list vars vals))
    else          (cons (list (first vars) (first vals))
                        (extend env (rest vars) (rest vals)))))

;; "invoke" in Queinnec
;; Note that our representation of functions here passes all args to the function
;; as a list as a single paramter
(defn qapply (funcarg args)
  (if (proc? funcarg)
    (funcarg args)
    (panic "not a function" funcarg)))

;; Takes a list of expressions and returns the corresponding list of values
(defn evlis (exps env)
  (if (list? exps)
    (cons (qeval (first exps) env)
          (evlis (rest exps) env))
    (list) ))



(defn make-function (funcargs funcbody env)
  ; TODO
  (eval (list 'fn (list funcargs) funcbody) env))

(defn qupdate! (id env value)
  (if (empty? env)
    (wrong "no such binding" id)
    (if (= (get (first env) 0) id)
      (begin (update-element! (first env) 1 value)
             value)
      (qupdate! id (rest env) value))))

(defn qbegin (exps env)
  (if (list? exps)
    (if (not (empty? (rest exps)))
      (begin
        (qeval (first exps) env)
        (qbegin (rest exps) env))
      (qeval (first exps) env))
    nil ; We return nil in the case of an empty begin
    ))

(let (env '((a 1)))
  (begin

    (defvtest "Atoms"
      (vt= (qeval '2 env) 2)
      (vt= (qeval "test" env) "test")
      (vt= (qeval \t env) \t)
      (vt= (qeval 'a env) 1)
      )

    (defvtest "Begin"
      (vt= (qeval '(begin 5 4) env)
           4))

  ))

(vt-start)
