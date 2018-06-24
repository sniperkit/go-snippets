(defn deco [a]
  (case a
    (\a \b \c \d \e \f) (char (+ (int a) 20))
    (\g \h \i \j \k \l \m) (char (+ (int a) 7))
    (\n \o \p \q \r \s \t) (char (- (int a) 7))
    (\u \v \w \x \y \z) (char (- (int a) 20))
    a))

(require '[clojure.string :as str])
(doseq [item (for [st (str/split-lines (slurp (first *command-line-args*)))]
  (map deco (seq st))
)] (println (str/join item)))
