(require '[clojure.string :as str])
(doseq [i (range 1 13)]
  (println
    (str/join ""
      (for [j (range 1 13)]
        (format "%4d" (* i j))))))
