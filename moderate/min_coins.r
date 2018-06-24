b <- c(0, 1, 2, 1, 2)
cat(sapply(as.integer(readLines(tail(commandArgs(), n=1))), function(s) {
  as.integer(s / 5) + b[s %% 5 + 1]
}), sep="\n")
