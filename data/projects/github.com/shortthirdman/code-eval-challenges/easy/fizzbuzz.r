fizz <- function(i, f, b, n) {
  if (i > n) {
    return("")
  } else if (i == n) {
    p <- ""
  } else {
    p <- " "
  }
  if (i %% f == 0 && i %% b == 0) {
    paste("FB", p, fizz(i+1, f, b, n), sep="")
  } else if (i %% f == 0) {
    paste("F", p, fizz(i+1, f, b, n), sep="")
  } else if (i %% b == 0) {
    paste("B", p, fizz(i+1, f, b, n), sep="")
  } else {
    paste(toString(i), p, fizz(i+1, f, b, n), sep="")
  }
}

fizzbuzz <- function(s) {
  fizz(1, s[1], s[2], s[3])
}

cat(sapply(lapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), as.integer), fizzbuzz), sep="\n")
