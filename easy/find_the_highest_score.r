cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " | ", fixed=TRUE), function(s) {
  c <- sapply(strsplit(s[1], " ")[[1]], as.integer)
  for (i in 2:length(s)) {
    d <- sapply(strsplit(s[i], " ")[[1]], as.integer)
    for (j in 1:length(d)) {
      if (d[j] > c[j]) {
        c[j] <- d[j]
      }
    }
  }
  paste(c, collapse=" ")
}), sep="\n")
