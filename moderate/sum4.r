cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), ","), function(st) {
  numzero <- function(i = 1, c = 4, z = 0) {
    if (c == 0) {
      if (z == 0) {
        1
      } else {
        0
      }
    } else if ((length(s) - i + 1 < c || z + c * s[i] > 0 || z + c * tail(s, n=1) < 0)) {
      0
    } else {
      numzero(i + 1, c - 1, z + s[i]) + numzero(i + 1, c, z)
    }
  }
  s <- sapply(st, as.integer)
  s <- s[order(s)]
  numzero()
}), sep="\n")
