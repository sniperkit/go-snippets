(require '[clojure.string :as str])

(defn wrd [a]
  (cond
    (= a "zero") 0
    (= a "one") 1
    (= a "two") 2
    (= a "three") 3
    (= a "four") 4
    (= a "five") 5
    (= a "six") 6
    (= a "seven") 7
    (= a "eight") 8
    (= a "nine") 9
    :else -1))

(doseq [item (for [st (remove empty? (str/split-lines (slurp (first *command-line-args*))))]
  (map wrd (str/split st #";"))
)] (println (str/join item)))
