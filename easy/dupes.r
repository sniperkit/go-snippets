cat(sapply(lapply(strsplit(readLines(tail(commandArgs(), n=1)), ","), unique), paste, collapse=","), sep="\n")
