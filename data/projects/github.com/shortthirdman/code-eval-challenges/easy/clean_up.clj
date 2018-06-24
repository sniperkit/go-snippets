(require '[clojure.string :as str])
(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (line-seq rdr)]
    (println (str/join " " (re-seq #"[a-z]+" (str/lower-case line))))))
