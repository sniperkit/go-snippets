cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(xs) {
  l <- strsplit(Reduce(function(x,y) {if (nchar(y) > nchar(x)) y else x}, xs, ""), "")[[1]]
  for (i in 1:length(l)) {
    l[i] <- paste(paste(rep("*", i - 1), collapse=""), l[i], sep="")
  }
  paste(l, collapse=" ")
}), sep="\n")
