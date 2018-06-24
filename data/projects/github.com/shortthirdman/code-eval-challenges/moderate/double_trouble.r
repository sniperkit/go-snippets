cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  n <- nchar(s) / 2
  r <- 1
  for (i in 1:n) {
    if ((substr(s, i, i) == "A" && substr(s, i+n, i+n) == "B") || (substr(s, i, i) == "B" && substr(s, i+n, i+n) == "A")) {
      r <- 0
      break
    } else if (substr(s, i, i) == "*" && substr(s, i+n, i+n) == "*") {
      r <- r * 2
    }
  }
  r
}), sep="\n")
