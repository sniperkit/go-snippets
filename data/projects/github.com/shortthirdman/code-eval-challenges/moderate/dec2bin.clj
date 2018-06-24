(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (remove empty? (line-seq rdr))]
    (println (Integer/toBinaryString (Integer/parseInt line)))))
