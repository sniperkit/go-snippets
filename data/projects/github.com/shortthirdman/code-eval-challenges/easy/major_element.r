for (line in readLines(tail(commandArgs(), n=1))) {
  cat(sapply(strsplit(line, ","), function(s) {
    a <- table(s)
    m <- which.max(a)
    if (a[m] > length(s)/2) {
      names(m)
    } else {
      "None"
    }
  }), sep="\n")
}
