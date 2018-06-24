cat(sapply(as.double(readLines(tail(commandArgs(), n=1))), function(t) {
  r <- 0
  if (t >= 2^31) {
    t <- t - 2^31
    r <- 1
  }
  s <- as.integer(t)
  while (s > 0) {
    s <- bitwAnd(s, s - 1)
    r <- r + 1
  }
  r %% 3
}), sep="\n")
