splitnum <- function(s) {
  r <- 0
  v <- 0
  o <- 1
  d <- 1
  for (i in 1:nchar(s[2])) {
    c <- substr(s[2], i, i)
    if (c == '+') {
      r <- r + v*o
      o <- 1
      v <- 0
    } else if (c == '-') {
      r <- r + v*o
      o <- -1
      v <- 0
    } else {
      v <- v*10 + as.integer(substr(s[1], d, d))
      d <- d + 1
    }
  }
  r + v*o
}

cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), splitnum), sep="\n")
