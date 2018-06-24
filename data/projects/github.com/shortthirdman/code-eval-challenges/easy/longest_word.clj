(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (apply max-key count (reverse (str/split st #" ")))
)] (println item))
