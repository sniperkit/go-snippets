for (line in readLines(tail(commandArgs(), n=1))) {
  cat(sapply(as.integer(line), function(s) {
    t <- c(1)
    if (s > 1) {
      for (i in 1:(s-1)) {
        r <- 1
        t <- append(t, 1)
        for (j in 1:i) {
          r <- (r * (i + 1 - j)) / j
          t <- append(t, r)
        }
      }
    }
    paste(t, collapse=" ")
  }), sep="\n")
}
