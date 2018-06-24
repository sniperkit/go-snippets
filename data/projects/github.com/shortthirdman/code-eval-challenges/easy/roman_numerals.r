ronum <- c(1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1)
rostr <- c("M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I")

cat(sapply(as.integer(readLines(tail(commandArgs(), n=1))), function(s) {
  i <- 1
  r <- c()
  while (s > 0) {
    if (s >= ronum[i]) {
      r <- append(r, rostr[i])
      s <- s - ronum[i]
    } else {
      i <- i + 1
    }
  }
  paste(r, collapse="")
}), sep="\n")
