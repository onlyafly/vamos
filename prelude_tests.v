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
               (= result true) (println "test passed")
               else (println "test failed"))
             (runtests othertests)))))

;;;;;;;;;; Tests

(deftest "Truthful values in 'if'"
  (= (list
       (if (quote true) 1 2)
       (if true 1 2))
     '(1 1)))

(runtests tests)
