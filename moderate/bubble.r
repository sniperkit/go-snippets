cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " | ", fixed=TRUE), function(s) {
  t <- sapply(strsplit(s[1], " ")[[1]], as.integer)
  if (length(t) > 1) {
  n <- min(as.integer(s[2]), length(t))
  for (i in 1:n) {
    for (j in 2:length(t)) {
      if (t[j-1] > t[j]) {
        tmp <- t[j-1]
        t[j-1] <- t[j]
        t[j] <- tmp
      }
    }
  }
  }
  paste(t, collapse=" ")
}), sep="\n")
