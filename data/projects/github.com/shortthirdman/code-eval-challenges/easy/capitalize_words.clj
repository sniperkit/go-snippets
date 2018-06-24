(require '[clojure.string :as str])
(defn capitalize [a]
  (str/join (cons (str/upper-case (first a)) (rest a))))
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (map capitalize (seq (str/split st #" ")))
)] (println (str/join " " item)))
