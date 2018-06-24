push <- function(s, x) {
  s <<- append(s, x)
}

pop <- function(s) {
  x <- s[length(s)]
  length(s) <<- length(s) - 1
  x
}

s <- c()
cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(ss) {
  s <<- c()
  for (i in 1:length(ss)) {
    s <<- push(s, ss[i])
  }
  h <- TRUE
  r <- c()
  for (i in 1:length(ss)) {
    x <- pop(s)
    if (h) {
      r <- append(r, x)
    }
    h <- !h
  }
  paste(r, collapse=" ")
}), sep="\n")
