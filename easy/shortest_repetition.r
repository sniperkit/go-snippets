cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), ""), function(s) {
  p <- 1
  if (length(s) > 1) {
    for (i in 2:length(s)) {
      if (s[i] != s[i-p]) {
        p <- i
      }
    }
  }
  p
}), sep="\n")
