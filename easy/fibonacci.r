fibo <- function(n, a=0, b=1) {
  if (n == 1) {
    b
  } else {
    fibo(n-1, b, a+b)
  }
}

cat(sapply(as.integer(readLines(tail(commandArgs(), n=1))), function(x) {
  if (x == 0) {
    0
  } else {
    fibo(x)
  }
}), sep="\n")
