n <- 0
m <- 0
for (s in readLines(tail(commandArgs(), n=1))) {
  if (n == 0) {
    n <- as.integer(s)
    l <- character(n)
  } else {
    for (i in 1:n) {
      if (l[i] == "") {
        l[i] <- s
        break
      } else if (nchar(s) > nchar(l[i])) {
        t <- s
        s <- l[i]
        l[i] <- t
      }
    }
  }
}
cat(l, sep="\n")
