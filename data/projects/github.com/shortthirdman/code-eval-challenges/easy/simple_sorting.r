cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(s) {
  paste(s[order(as.double(s))], collapse=" ")
}), sep="\n")
