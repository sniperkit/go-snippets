cat(sapply(as.integer(readLines(tail(commandArgs(), n=1))), function(s) {
  r <- 0
  while (s > 0) {
    s <- bitwAnd(s, s - 1)
    r <- r + 1
  }
  r
}), sep="\n")
