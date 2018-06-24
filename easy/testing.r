cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  t <- strsplit(s, " | ", fixed=TRUE)[[1]]
  r <- 0
  for (i in 1:nchar(t[1])) {
    if (substr(t[1], i, i) != substr(t[2], i, i)) {
      r <- r + 1
    }
  }
  if (r == 0) {
    "Done"
  } else if (r <= 2) {
    "Low"
  } else if (r <= 4) {
    "Medium"
  } else if (r <= 6) {
    "High"
  } else {
    "Critical"
  }
}), sep="\n")
