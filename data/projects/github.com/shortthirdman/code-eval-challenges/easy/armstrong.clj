(defn armstrong [a]
  (def e (Math/ceil (Math/log10 a)))
  (def i a)
  (def r a)
  (while (> i 0)
    (when true
      (def r (- r (Math/pow (mod i 10) e)))
      (def i (quot i 10))))
  (if (== r 0)
    "True"
    "False"))

(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (line-seq rdr)]
    (println (armstrong (Integer/parseInt line)))))
