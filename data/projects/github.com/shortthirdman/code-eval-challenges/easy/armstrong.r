cat(sapply(as.integer(readLines(tail(commandArgs(), n=1))), function(s) {
  if (s < 10) { return("True") }
  e <- as.integer(log10(s)) + 1
  r <- s
  while (s > 0) {
    r <- r - (s %% 10)^e
    s <- as.integer(s / 10)
  }
  if (r == 0) { return("True") }
  "False"
}), sep="\n")
