(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (#(when true
      (def n (Integer/parseInt (first %)))
      (def m (Integer/parseInt (second %)))
      (-  n (* m (quot n m))))
(str/split st #",")))] (prn item))
