cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " : ", fixed=TRUE), function(s) {
  paste(Reduce(function(x,y) {
    tmp <- x[y[1]+1]
    x[y[1]+1] <- x[y[2]+1]
    x[y[2]+1] <- tmp
    x
  }, Map(function(x) {
    Map(as.integer, strsplit(x, "-"))
  }, strsplit(s[2], ", ", fixed=TRUE))[[1]], strsplit(s[1], " ")[[1]]), collapse=" ")
}), sep="\n")
