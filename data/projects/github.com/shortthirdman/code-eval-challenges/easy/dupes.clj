(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (distinct (str/split st #","))
)] (println (apply str (interpose "," item))))
