cat(sapply(lapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), as.integer), function(s) {
  j <- TRUE
  if (s[1] > 1) {
    u <- c()
    for (i in 3:length(s)) {
      x <- abs(s[i] - s[i-1])
      if (x >= s[1] || x == 0 || (x-1) %in% u) {
        j <- FALSE
        break
      }
      u <- append(u, x-1)
    }
  }
  if (j == TRUE) {
    "Jolly"
  } else {
    "Not jolly"
  }
}), sep="\n")
