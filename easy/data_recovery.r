for (line in readLines(tail(commandArgs(), n=1))) {
  cat(sapply(strsplit(line, ";"), function(s) {
    t <- strsplit(s[1], " ")[[1]]
    u <- sapply(strsplit(s[2], " ")[[1]], as.integer)
    r <- c()
    for (i in 1:length(u)) {
      r[u[i]] <- t[i]
    }
    for (i in 1:length(t)) {
      if (is.na(r[i])) {
        r[i] <- t[length(t)]
        break
      }
    }
    paste(r, collapse=" ")
  }), sep="\n")
}
