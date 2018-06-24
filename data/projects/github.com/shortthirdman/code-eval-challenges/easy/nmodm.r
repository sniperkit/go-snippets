cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), ","), function(s) {
  n <- as.integer(s[1])
  m <- as.integer(s[2])
  n - as.integer(n / m) * m
}), sep="\n")
