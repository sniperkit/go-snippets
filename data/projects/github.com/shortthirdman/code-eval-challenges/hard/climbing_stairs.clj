(def cli-seq
  (lazy-cat [0N 1N] (map + (rest cli-seq) cli-seq)))
(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (remove empty? (line-seq rdr))]
    (println (format "%d" (biginteger (nth cli-seq (inc (Integer/parseInt line))))))))
