(require '[clojure.string :as str])
(defn reverse-cmp [a b]
  (compare b a))
(doseq [item (for [st (remove empty? (str/split-lines (slurp (first *command-line-args*))))]
  (sort reverse-cmp (str/split st #" "))
)] (println (str/join " " item)))
