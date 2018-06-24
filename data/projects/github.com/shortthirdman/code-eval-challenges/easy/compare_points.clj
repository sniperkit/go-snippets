(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (when true
    (def p (map #(Integer/parseInt (str %)) (str/split st #" ")))
    (cond
      (= (first p) (nth p 2))
        (if (= (second p) (nth p 3))
          "here"
          (if (< (second p) (nth p 3))
            "N"
            "S"))
      (= (second p) (nth p 3))
        (if (< (first p) (nth p 2))
          "E"
          "W")
      :else
        (if (< (first p) (nth p 2))
          (if (< (second p) (nth p 3))
            "NE"
            "SE")
          (if (< (second p) (nth p 3))
            "NW"
            "SW"))))
)] (println item))
