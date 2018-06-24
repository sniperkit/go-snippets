(defn colu [s]
  (if (= s 0)
    ""
    (str (colu (quot (- s 1) 26)) (char (+ (int \A) (mod (- s 1) 26))))))
(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (remove empty? (line-seq rdr))]
    (println (colu (Integer/parseInt line)))))
