(def fib-seq
  (lazy-cat [0 1] (map + (rest fib-seq) fib-seq)))
(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (line-seq rdr)]
    (prn (nth fib-seq (Integer/parseInt line)))))
