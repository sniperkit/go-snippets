cat(sapply(lapply(strsplit(readLines(tail(commandArgs(), n=1)), "[,;]"), function(s) s[duplicated(s)]), paste, collapse=","), sep="\n")
