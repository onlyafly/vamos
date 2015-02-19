;;;;;;;;;;;;;;;;;;;;;;;;
; vtest test framework ;
;;;;;;;;;;;;;;;;;;;;;;;;

(def _vtest_tests '())

(defn _vtest_runtests (tests)
  (cond
    (= tests '()) nil
    else (let (test (first tests)
               othertests (rest tests)
               testname (first test)
               testfn (first (rest test))
               result (testfn))
           (begin
             (cond
               (= result true) (println ".")
               else (println (concat "TEST FAILED: " testname)))
             (_vtest_runtests othertests)))))

;;;;;;;;;; External API

;; (defvtest "Sample Test"
;;   pred1 pred2 predn...)
;; =>
;; (update! _vtest_tests
;;          (cons (list "Sample Test" (fn () (begin pred1 pred2 predn...)))
;;                _vtest_tests))
;;
(defmacro defvtest (name &rest preds)
  (list 'update! '_vtest_tests
    (list 'cons
      (list 'list name
        (list 'fn '()
          (cons 'begin preds)))
      '_vtest_tests)))

(defn vt= (actual expected)
  (if (= actual expected)
    true
    (begin
      (println (concat "TEST FAILED"))
      false)))

(defn vt-start ()
  (println "Running vtest tests...")
  (_vtest_runtests _vtest_tests)
  (println "Tests complete..."))
