cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), ";"), function(s) {
  t <- strsplit(s[2], ",")[[1]]
  t[duplicated(t)][1]
}), sep="\n")
