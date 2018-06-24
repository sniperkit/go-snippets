cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(xs) {
  if (length(xs) < 2) {
    ""
  } else {
    xs[length(xs)-1]
  }
}), sep="\n")
