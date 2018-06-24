cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(xs) {
  Reduce(function(x,y) {if (nchar(y) > nchar(x)) y else x}, xs, "")
}), sep="\n")
