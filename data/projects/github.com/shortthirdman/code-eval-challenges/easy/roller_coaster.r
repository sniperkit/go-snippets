cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  u <- FALSE
  paste(lapply(strsplit(s, "")[[1]], function(x) {
    if (x %in% LETTERS || x %in% letters) {
      u <<- !u
      if (u) {
        toupper(x)
      } else {
        tolower(x)
      }
    } else {
      x
    }
  }), collapse="")
}), sep="\n")
