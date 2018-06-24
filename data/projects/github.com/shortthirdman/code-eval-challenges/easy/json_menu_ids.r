cat(Filter(function(x) !is.na(x), sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  if (nchar(s) == 0) {
    return(NA)
  }
  m <- gregexpr("\"id\": \\d+,", s)[[1]]
  if (attr(m, "match.length") == -1) {
    return("0")
  }
  r <- 0
  for (i in m) {
    r <- r + as.integer(gsub("\\D+", "", substr(s, i+5, i+11)))
  }
  r
})), sep="\n")
