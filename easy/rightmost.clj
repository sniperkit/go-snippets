(require '[clojure.string :as str])
(doseq [item (for [st (remove empty? (str/split-lines (slurp (first *command-line-args*))))]
  (#(.lastIndexOf (first %) (second %))
(str/split st #",")))] (prn item))
