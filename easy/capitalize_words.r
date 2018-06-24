cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(s) {
  paste(sapply(s, function(s) {
    paste(toupper(substr(s, 1, 1)), substr(s, 2, nchar(s)), sep="")
  }), collapse=" ")
}), sep="\n")
