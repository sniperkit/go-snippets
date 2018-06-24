(defn in? [seq elm]
  (some #(= elm %) seq))
(defn happy [a]
  (if (< a 10)
    (* a a)
    (+ (* (mod a 10) (mod a 10)) (happy (quot a 10)))))
(defn ishappy [a l]
  (if (= a 1)
    1
    (if (in? l a)
      0
      (ishappy (happy a) (conj l a)))))

(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (line-seq rdr)]
    (prn (ishappy (Integer/parseInt line) []))))
