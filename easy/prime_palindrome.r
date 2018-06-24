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

isPalindrome <- function(s) {
  return(s == sapply(lapply(strsplit(s, ""), rev), paste, collapse=""))
}

maxPrime <- 0
nextPrime <- primeSeq()
j <- nextPrime()
while (j < 1000) {
  if (isPalindrome(toString(j))) {
    maxPrime <- j
  }
  j <- nextPrime()
}
cat(maxPrime, sep="\n")
