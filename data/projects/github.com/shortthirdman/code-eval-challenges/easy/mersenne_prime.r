primes <- c(2, 3)
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

mersenne <- function(x, n) { 2^x - 1 < n }
nextPrime <- primeSeq()
m <- 5
cat(sapply(as.integer(readLines(tail(commandArgs(), n=1))), function(n) {
  while (2^m - 1 < n) {
    isPrime(m)
    m <<- m + 2
  }
  paste(2^primes[mersenne(primes, n)] - 1, collapse=", ")
}), sep="\n")
