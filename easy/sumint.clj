(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (prn (reduce + (map #(Integer/parseInt %) (line-seq rdr)))))
