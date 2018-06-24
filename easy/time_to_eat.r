cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(s) {
  paste(sort(s, decreasing=TRUE), collapse=" ")
}), sep="\n")
