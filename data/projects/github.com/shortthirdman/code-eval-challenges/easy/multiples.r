multiples <- function(s) {
  c <- s[1] - floor(s[1] * 2^(-log2(s[2]))) * 2^log2(s[2])
  if (c > 0) {
    s[1] - c + s[2]
  } else {
    s[1]
  }
}

cat(sapply(lapply(strsplit(readLines(tail(commandArgs(), n=1)), ","), as.integer), multiples), sep="\n")
