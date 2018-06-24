numf <- function(s) {
  sapply(strsplit(gsub("[^ 0123456789]", "", s), " ")[[1]], as.integer)
}

for (line in readLines(tail(commandArgs(), n=1))) {
  cat(sapply(lapply(line, numf), function(s) {
    s <- sort.int(s)
    m <- s[1]
    for (i in 2:length(s)) {
      s[i] <- s[i] - m
      m <- m + s[i]
    }
    paste(s, collapse=",")
  }), sep="\n")
}
