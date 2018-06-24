cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  Reduce("+", lapply(strsplit(s, "")[[1]], as.integer))
}), sep="\n")
