contains_all <- function(a, b) {
  if (b == "") {
    return(TRUE)
  }
  c <- sub(substr(b, 1, 1), "", a)
  if (c == a) {
    return(FALSE)
  }
  contains_all(c, substr(b, 2, nchar(b)))
}

cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  t <- strsplit(s, " | ", fixed=TRUE)[[1]]
  u <- strsplit(t[1], " ")[[1]]
  r <- ""
  for (i in u) {
    if (contains_all(i, t[2])) {
      if (r == "") {
        r <- i
      } else {
        r <- paste(r, i)
      }
    }
  }
  if (r == "") {
    "False"
  } else {
    paste(r, collapse=" ")
  }
}), sep="\n")
