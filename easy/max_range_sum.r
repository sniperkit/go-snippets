for (line in readLines(tail(commandArgs(), n=1))) {
  cat(sapply(strsplit(line, ";"), function(s) {
    n <- as.integer(s[1])
    d <- sapply(strsplit(s[2], " ")[[1]], as.integer)
    c <- sum(d[1:n])
    m <- max(0, c)
    for (i in (n+1):length(d)) {
      c <- c - d[i-n] + d[i]
      m <- max(m, c)
    }
    m
  }), sep="\n")
}
