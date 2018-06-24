cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(s) {
  paste(apply(t(matrix(s, sqrt(length(s)))), 2, rev), collapse=" ")
}), sep="\n")
