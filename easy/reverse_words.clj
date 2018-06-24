(require '[clojure.string :as str])
(doseq [item (for [st (remove empty? (str/split-lines (slurp (first *command-line-args*))))]
  (reverse (str/split st #" "))
)] (println (str/join " " item)))
