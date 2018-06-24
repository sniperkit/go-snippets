cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " | ", fixed=TRUE), function(s) {
  t <- sapply(strsplit(s[1], " ")[[1]], as.integer)
  n <- as.integer(s[2])
  l <- 2
  for (i in 1:n) {
    p <- l
    l <- 0
    for (j in p:length(t)) {
      if (t[j-1] > t[j]) {
        tmp <- t[j-1]
        t[j-1] <- t[j]
        t[j] <- tmp
        l <- if(j > 2) j - 1 else 3
        break
      }
    }
    if (l == 0) { break }
  }
  paste(t, collapse=" ")
}), sep="\n")
