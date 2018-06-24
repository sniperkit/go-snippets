cat(Filter(function(x) !is.na(x), sapply(strsplit(readLines(tail(commandArgs(), n=1)), ","), function(xs) {
  if (length(xs) < 2) {
    NA
  } else if (grepl(paste(xs[2], "$", sep=""), xs[1])) {
    "1"
  } else {
    "0"
  }
})), sep="\n")
