m <- c(50, 25, 10, 5, 1)
alter <- function(n, p) {
  if (n == 0) {
    return(1)
  }
  while (m[p] > n) {
    p <- p + 1
  }
  if (m[p] == 1) {
    return(1)
  }
  alter(n - m[p], p) + alter(n, p + 1)
}

cat(sapply(as.integer(readLines(tail(commandArgs(), n=1))), function(s) alter(s, 1)), sep="\n")
