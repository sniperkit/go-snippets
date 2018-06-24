cat(sapply(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), as.integer), function(s) {
  n <- s[1]
  s <- sort.int(s[-1])
  u <- s[n/2 + 1]
  toString(Reduce(function(a,b) {a + abs(b - u)}, s, 0))
}), sep="\n")
