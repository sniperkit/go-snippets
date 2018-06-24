cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  t <- strsplit(s, " ")[[1]]
  r <- ""
  for (i in 1:nchar(t[1])) {
    if (substr(t[2], i, i) == '1') {
      r <- paste(r, toupper(substr(t[1], i, i)), sep='')
    } else {
      r <- paste(r, tolower(substr(t[1], i, i)), sep='')
    }
  }
  r
}), sep="\n")
