cat(Filter(function(x) !is.na(x), sapply(lapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), rev), function(s) {
  if (length(s) > 0) {
    paste(s, collapse=" ")
  } else {
    NA
  }
})), sep="\n")
