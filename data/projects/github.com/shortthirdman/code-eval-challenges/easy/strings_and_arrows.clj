(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (+ (count (re-seq #"(?=(>>-->))" st)) (count (re-seq #"(?=(<--<<))" st)))
)] (prn item))
