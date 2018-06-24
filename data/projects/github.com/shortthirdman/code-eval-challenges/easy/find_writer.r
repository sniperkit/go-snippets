cat(na.omit(sapply(strsplit(readLines(tail(commandArgs(), n=1)), "| ", fixed=TRUE), function(s) {
  if (length(s) > 1) {
    paste(lapply(strsplit(s[2], " ")[[1]], function(t) {
      substr(s[1], as.integer(t), as.integer(t))
    }), collapse="")
  } else {
    NA
  }
})), sep="\n")
