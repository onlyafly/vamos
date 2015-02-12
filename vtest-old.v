;;;;;;;;;; Test framework

(def tests '())

(defmacro deftest (name pred)
  (list 'update! 'tests
    (list 'cons
      (list 'list name
        (list 'fn '() pred))
      'tests)))

(defn runtests (tests)
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
             (runtests othertests)))))
