(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (reduce + (map #(Integer/parseInt (str %)) (seq st)))
)] (prn item))
