paren <- function(xs, ys) {
  if (length(ys) == 0 || is.na(ys[1])) {
    return(length(xs) == 0 || is.na(xs[1]))
  }
  if (length(xs) > 0) {
    if ((ys[1] == ")" && xs[1] == "(") || (ys[1] == "]" && xs[1] == "[") || (ys[1] == "}" && xs[1] == "{")) {
      return(paren(xs[2:length(xs)], ys[2:length(ys)]))
    }
  }
  if (ys[1] %in% ")]}") {
    return(FALSE)
  }
  paren(c(ys[1], xs, recursive=TRUE), ys[2:length(ys)])
}

cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  if (paren(vector(), strsplit(s, "")[[1]])) {
    "True"
  } else {
    "False"
  }
}), sep="\n")
