(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (remove empty? (line-seq rdr))]
    (prn (reduce + (map #(Integer/parseInt (last %)) (re-seq #"\"id\": (\d+)," line))))))
