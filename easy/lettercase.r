cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  l <- sum(attr(gregexpr("[a-z]", s)[[1]], "match.length"))
  u <- sum(attr(gregexpr("[A-Z]", s)[[1]], "match.length"))
  if (l == -1) { l <- 0 }
  if (u == -1) { u <- 0 }
  if (l+u > 0) {
    sprintf("lowercase: %.2f uppercase: %.2f", 100*l/(l+u), 100*u/(l+u))
  } else {
    "lowercase: 0.00 uppercase: 0.00"
  }
}), sep="\n")
