(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (line-seq rdr)]
    (prn (- 1 (bit-and (Integer/parseInt line) 1)))))
