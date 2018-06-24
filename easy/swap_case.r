cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  paste(lapply(strsplit(s, "")[[1]], function(x) {
    if (x %in% LETTERS) {
      tolower(x)
    } else if (x %in% letters) {
      toupper(x)
    } else {
      x
    }
  }), collapse="")
}), sep="\n")
