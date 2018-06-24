(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (#(= (bit-test (Integer/parseInt (first %))
                 (dec (Integer/parseInt (second %))))
       (bit-test (Integer/parseInt (first %))
                 (dec (Integer/parseInt (nth % 2)))))
(str/split st #",")))] (prn item))
