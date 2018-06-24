cat(sapply(lapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), as.integer), function(s) {
  l <- s[1]
  d <- s[2]
  n <- s[3]
  count <- 0
  t <- 0
  tx <- 6 - d
  i <- 6
  while (i <= l - 6) {
    if (i > tx - d) {
      i <- tx
      if (t == n) {
        tx <- l - 6 + d
      } else {
        tx <- s[t+4]
        t <- t + 1
      }
    } else {
      count <- count + 1
    }
    i <- i + d
  }
  count
}), sep="\n")
