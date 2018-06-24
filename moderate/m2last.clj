(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (#(when true
      (def l (Integer/parseInt (last %)))
      (if (< l (count %))
        (nth % (- (count %) (inc l)))
        nil))
(str/split st #" ")))]
  (if (not (empty? item))
    (println item)))
