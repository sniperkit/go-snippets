(def b [0 1 2 1 2])
(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (line-seq rdr)]
    (when true
      (def n (Integer/parseInt line))
      (prn (+ (quot n 5) (nth b (mod n 5)))))))
