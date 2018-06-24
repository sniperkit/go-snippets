(defn age [a]
  (cond
    (or (< a 0) (> a 100)) "This program is for humans"
    (< a 3) "Still in Mama's arms"
    (< a 5) "Preschool Maniac"
    (< a 12) "Elementary school"
    (< a 15) "Middle school"
    (< a 19) "High school"
    (< a 23) "College"
    (< a 66) "Working for the man"
    :else "The Golden Years"))

(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (line-seq rdr)]
    (println (age (Integer/parseInt line)))))
