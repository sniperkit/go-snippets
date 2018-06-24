xz <- function(n, a) {
  while (a > 0) {
    if (a %% 2 == 0) {
      if (n > 0) {
        n <- n - 1
      } else {
        return(FALSE)
      }
    }
    a <- a %/% 2
  }
  n == 0
}

cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  t <- sapply(strsplit(s, " ")[[1]], as.integer)
  r <- 0
  for (i in 2^(t[1]-1):t[2]) {
    if (xz(t[1], i)) {
      r <- r + 1
    }
  }
  r
}), sep="\n")
