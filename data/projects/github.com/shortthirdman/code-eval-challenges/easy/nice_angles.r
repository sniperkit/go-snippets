cat(sapply(as.numeric(readLines(tail(commandArgs(), n=1))), function(s) {
  a <- as.integer(s)
  r <- paste(toString(a), ".", sep="")
  s <- (s - a) * 60
  a <- as.integer(s)
  if (a < 10) { f <- "0" } else { f <- "" }
  r <- paste(r, f, toString(a), "'", sep="")
  s <- (s - a) * 60
  a <- as.integer(s)
  if (a < 10) { f <- "0" } else { f <- "" }
  r <- paste(r, f, toString(a), "\"", sep="")
  r
}), sep="\n")
