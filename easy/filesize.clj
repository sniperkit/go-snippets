(import '[java.io File])
(prn (.length (File. (first *command-line-args*))))
