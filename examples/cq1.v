;; cq1
;; From Chapter 1 of "Lisp in Small Pieces" by Christian Queinnec
;; Updated 2016-02-09
;;
;; Run via: vamos examples/cq1.v
#|
Features of the CQ1 language:

Dynamic binding function:

((let (y 2)
   (let (y 1)
     (dynamic-fn (x) (list x y))))
 3)
=> (3 2)

Lexical binding function:

((let (y 2)
   (let (y 1)
     (fn (x) (list x y))))
 3)
=> (3 1)

((((fn (a)
     (fn (a)
       (fn (b) a)))
   1)
  2)
 3)

|#

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
  ;;DEBUG (println "qeval:" e "::" env)
  (if (atom? e)

    ;; Handle atoms
    (cond
      (symbol? e) (lookup e env)
      (or (number? e) (string? e) (boolean? e) (char? e)) e
      else (wrong "cannot evaluate" e))

    ;; Handle non-atoms
    (let (proc (first e))
      (cond
        (= proc 'apply)      (invoke (qeval (get e 1) env)
                                     (qeval (get e 2) env)
                                     env)
        (= proc 'quote)      (get e 1)
        (= proc 'if)         (if (qeval (get e 1) env)
                               (qeval (get e 2) env)
                               (qeval (get e 3) env))
        (= proc 'begin)      (qbegin (rest e) env)
        (= proc 'update!)    (qupdate! (get e 1) env (qeval (get e 2) env))

        ;; Dynamically-scoped function
        ;; Syntax: (dynamic-fn (param1 ...) body ...)
        (= proc 'dynamic-fn) (let (params (get e 1)
                                   body (rest (rest e)))
                               (make-dynamic-function params body env))

        ;; Lexically-scoped function
        ;; Syntax: (fn (param1 ...) body ...)
        (= proc 'fn)         (let (params (get e 1)
                                   body (rest (rest e)))
                               (make-function params body env))

        else                 (invoke (qeval (first e) env)
                                     (evlis (rest e) env)
                                     env)
        ))))

;; A CQ1-function is represented as a Vamos-function, where the CQ1-function's arguments are passed as a list
(defn make-dynamic-function (funcparams funcbody lexical.env)
  (fn (args dynamic.env)
    (qbegin funcbody (extend dynamic.env funcparams args))))

;; A CQ1-function is represented as a Vamos-function, where the CQ1-function's arguments are passed as a list
(defn make-function (funcparams funcbody lexical.env)
  (fn (args dynamic.env)
    (qbegin funcbody (extend lexical.env funcparams args))))

(defn make-primitive (primname primfunc min-arity max-arity)
  (fn (args dynamic.env)
    (if (and (<= min-arity (len args))
             (>= max-arity (len args)))
      (apply primfunc args)
      (wrong "Incorrect arity" (list primname args)))))

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

;; Note that our representation of functions here passes all args to the function
;; as a list as a single paramter
(defn invoke (funcarg args dynamic.env)
  (if (proc? funcarg)
    (funcarg args dynamic.env)
    (panic "not a function" funcarg)))

;; Takes a list of expressions and returns the corresponding list of values
(defn evlis (exps env)
  ;;DEBUG (println "evlis:" exps "::" env)
  (if (empty? exps)
    (list)
    (cons (qeval (first exps) env)
          (evlis (rest exps) env))))

(defn qupdate! (id env value)
  (if (empty? env)
    (wrong "no such binding" id)
    (if (= (get (first env) 0) id)
      (begin (update-element! (first env) 1 value)
             value)
      (qupdate! id (rest env) value))))

(defn qbegin (exps env)
  ;;DEBUG (println "qbegin:" exps "::" env)
  (if (list? exps)
    (if (not (empty? (rest exps)))
      (begin
        (qeval (first exps) env)
        (qbegin (rest exps) env))
      (qeval (first exps) env))
    nil ; We return nil in the case of an empty begin
    ))

(def env.init '())

(def env.global env.init)

;; value is optional
(defmacro definitial (name value)
  (list 'begin (list 'update! 'env.global (list 'cons (list 'list (list 'quote name) value)
                                                      'env.global))
               (list 'quote name)))

(defmacro defprimitive (name f min-arity max-arity)
  (list 'definitial
        name
        (list 'make-primitive (list 'quote name) f min-arity max-arity)))

(definitial foo nil)
(definitial bar nil)

(defprimitive cons cons 2 2)
(defprimitive first first 1 1)
(defprimitive update-element! update-element! 3 3)
(defprimitive + + 2 2)
(defprimitive = = 2 2)
(defprimitive < < 2 2)
(defprimitive list list 0 1000)

(let (env '((a 1))
      genv env.global)
  (begin

    (defvtest "Atoms"
      (vt= (qeval '2 env) 2)
      (vt= (qeval "test" env) "test")
      (vt= (qeval \t env) \t)
      (vt= (qeval 'a env) 1)
      )

    (defvtest "Quote"
      (vt= (qeval '(quote (1 2)) genv)
           '(1 2)
           ))


    (defvtest "Apply"
      (vt= (qeval '(apply + (quote (1 2))) genv)
           3))

    (defvtest "Begin"
      (vt= (qeval '(begin 5 4) env)
           4))

    (defvtest "Demonstrate lexical binding"
      (vt= (qeval '(((fn (a) (fn (b) a)) 1) 2)
                  '((a 3)))
           1))

    (defvtest "Demonstrate dynamic binding"
      (vt= (qeval '(((fn (a) (dynamic-fn (b) a)) 1) 2)
                  '((a 3)))
           3))

  ))

(vt-start)

(defn toplevel ()
  (println (qeval (read-string (read-line)) env.global))
  (toplevel))
(toplevel)
