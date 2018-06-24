cat(Filter(function(x) !is.na(x), sapply(strsplit(readLines(tail(commandArgs(), n=1)), ","), function(s) {
  if (length(s) > 1) {
    a <- tail(gregexpr(s[2], s[1])[[1]], n=1)
    if (a == -1) {
      -1
    } else {
      a-1
    }
  } else {
    NA
  }
})), sep="\n")
