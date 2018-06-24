(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (str/replace st #"(.)\1+" "$1")
)] (println item))
