cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " | ", fixed=TRUE), function(s) {
  a <- as.list(strsplit(s[1], " ")[[1]])
  m <- as.integer(s[2])
  while(length(a) > 1) {
    n <- m %% length(a)
    a[[ifelse(n > 0, n, length(a))]] <- NULL
  }
  a[[1]]
}), sep="\n")
