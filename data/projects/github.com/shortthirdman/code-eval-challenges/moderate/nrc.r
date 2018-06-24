cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), ""), function(s) {
  t <- table(s)
  for (i in s) {
    if (t[i] == 1) {
      return(i)
    }
  }
  ""
}), sep="\n")
