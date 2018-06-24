cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), ", "), function(xs) {
  paste(Filter(function(c) !(c %in% strsplit(xs[2], "")[[1]]), strsplit(xs[1], "")[[1]]), collapse="")
}), sep="\n")
