cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  xs <- strsplit(s, ",")[[1]]
  x <- as.integer(xs[1])
  p1 <- 2^(as.integer(xs[2])-1)
  p2 <- 2^(as.integer(xs[3])-1)
  if ((x %% (p1+p1) >= p1) == (x %% (p2+p2) >= p2)) {
    "true"
  } else {
    "false"
  }
}), sep="\n")
