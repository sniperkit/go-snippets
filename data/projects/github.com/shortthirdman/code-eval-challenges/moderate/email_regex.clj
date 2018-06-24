(with-open [rdr (clojure.java.io/reader (first *command-line-args*))]
  (doseq [line (line-seq rdr)]
    (println
      (if (re-find #"^((\"[0-9A-Za-z@.+-=]+\")|([0-9A-Za-z.+-=]+))@\w*(\w+\.)+\w{2,4}$" line)
        "true"
        "false"))))
