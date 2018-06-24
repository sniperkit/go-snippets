(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (line-seq rdr)]
    (prn (- (reduce min (map count (re-seq #"X\.*Y" line))) 2))))
