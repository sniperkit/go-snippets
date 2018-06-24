cat(sapply(tolower(readLines(tail(commandArgs(), n=1))), function(s) {
  xs <- strsplit(s, "")[[1]]
  ct <- rep.int(0, 26)
  for (i in 1:length(xs)) {
    if (xs[i] >= 'a' && xs[i] <= 'z') {
      ct[utf8ToInt(xs[i]) - 96] <- ct[utf8ToInt(xs[i]) - 96] + 1
    }
  }
  ct <- sort.int(ct)
  i <- 26
  b <- 0
  while (i > 0 && ct[i] > 0) {
    b <- b + i * ct[i]
    i <- i - 1
  }
  b
}), sep="\n")
