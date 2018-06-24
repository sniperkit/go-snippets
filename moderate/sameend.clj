(require '[clojure.string :as str])
(doseq [item (for [st (remove empty? (str/split-lines (slurp (first *command-line-args*))))]
  (#(.endsWith (first %) (second %))
(str/split st #",")))] (prn (if item 1 0)))
