a <- strtoi(charToRaw("A"), 16L)

cat(sapply(as.integer(readLines(tail(commandArgs(), n=1))), function(s) {
  s <- s - 1
  r <- rawToChar(as.raw(s %% 26 + a))
  while (s >= 26) {
    s <- s %/% 26 - 1
    r = paste(rawToChar(as.raw(s %% 26 + a)), r, sep="")
  }
  r
}), sep="\n")
