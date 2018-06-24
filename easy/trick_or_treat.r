numi <- function(s) {
  strsplit(gsub("[^,0123456789]", "", s), ",")[[1]]
}
cat(sapply(lapply(readLines(tail(commandArgs(), n=1)), numi), function(s) {
  t <- as.integer(s)
  ifelse(t[1] + t[2] + t[3] == 0, 0,
         ((t[1] * 3 + t[2] * 4 + t[3] * 5) * t[4]) %/% (t[1] + t[2] + t[3]))
}), sep="\n")
