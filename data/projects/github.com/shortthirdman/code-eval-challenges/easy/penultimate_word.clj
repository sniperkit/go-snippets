(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (last (butlast (str/split st #" ")))
)] (println item))
