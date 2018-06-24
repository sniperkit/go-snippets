cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  f <- TRUE
  w <- TRUE
  paste(lapply(strsplit(s, "")[[1]], function(x) {
    if (x %in% LETTERS) {
      f <<- FALSE
      if (!w) {
        w <<- TRUE
        paste(" ", tolower(x), sep="")
      } else {
        tolower(x)
      }
    } else if (x %in% letters) {
      f <<- FALSE
      if (!w) {
        w <<- TRUE
        paste(" ", x, sep="")
      } else {
        x
      }
    } else {
      if (!f) {
        w <<- FALSE
      }
      ""
    }
  }), collapse="")
}), sep="\n")
