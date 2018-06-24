juggle <- function(n, s) {
  if (is.na(s[1])) {
    n
  } else if (s[1] == 1) {
    juggle(n * 2^s[2], s[3:length(s)])
  } else {
    juggle((n+1) * 2^s[2] - 1, s[3:length(s)])
  }
}

cat(sapply(lapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), nchar), function(s) {
  toString(juggle(0, s))
}), sep="\n")
