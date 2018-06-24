primes <- c(2, 3, 5, 7, 11, 13)
isPrime <- function(a) {
  i <- 1
  while (a %% primes[i] > 0) {
    if (primes[i] * primes[i] > a) {
      primes <<- append(primes, a)
      return(TRUE)
    }
    i <- i + 1
  }
  FALSE
}

primeSeq <- function() {
  p0 <- 1
  i <- 0
  return(function() {
    if (p0 <= length(primes)) {
      i <<- primes[p0]
      p0 <<- p0 + 1
      return(i)
    }
    i <<- i + 2
    while (!isPrime(i)) {
      i <<- i + 2
    }
    p0 <<- p0 + 1
    return(i)
  })
}

primeless <- function(n) {
  r <- c()
  nextPrime <- primeSeq()
  i <- nextPrime()
  while (i < n) {
    r <- append(r, i)
    i <- nextPrime()
  }
  paste(r, collapse=",")
}

for (line in readLines(tail(commandArgs(), n=1))) {
  cat(sapply(as.integer(line), primeless), sep="\n")
}
