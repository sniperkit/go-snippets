s <- .Platform$endian
cat(paste(toupper(substring(s, 1, 1)), substring(s, 2), "Endian", sep = ""), sep="\n")
