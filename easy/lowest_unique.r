cat(sapply(lapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), as.integer), function(s) {
  l <- 0
  lx <- 0
  for (i in 1:length(s)) {
    if (length(which(s == s[i])) == 1 && (lx == 0 || l > s[i])) {
      l <- s[i]
      lx <- i
    }
  }
  lx
}), sep="\n")
