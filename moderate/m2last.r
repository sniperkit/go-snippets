cat(Filter(function(x) !is.na(x), sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(s) {
  i <- as.integer(tail(s, n=1))
  if (i >= length(s)) {
    return(NA)
  }
  s[length(s) - i]
})), sep="\n")
