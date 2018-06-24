cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  t <- sapply(strsplit(s, "")[[1]], as.integer)
  u <- rep(0, length(t))
  for (i in t) {
    v <- i + 1
    if (v <= length(t)) {
      u[v] <- u[v] + 1
    } else {
      break
    }
  }
  if (all(t == u)) {
    "1"
  } else {
    "0"
  }
}), sep="\n")
