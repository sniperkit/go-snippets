cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(s) {
  paste(Reduce(function(x,y) {
    if (x[1] == 0) {
      c(1, y)
    } else {
      if (x[length(x)] == y) {
        x[length(x)-1] <- as.integer(x[length(x)-1]) + 1
        x
      } else {
        append(x, c(1, y))
      }
    }
  }, s, c(0, "")), collapse=" ")
}), sep="\n")
