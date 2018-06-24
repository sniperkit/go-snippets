p <- c("8", "30", "20", "8", "20", NA, "30")
names(p) <- c("3", "8", "10", "20", "29", "30", "52")

cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(s) {
  while (!is.na(s[1])) {
    if (s[1] == s[2]) {
      break
    }
    i <- s[2]
    while (!is.na(i)) {
      if (s[1] == i) {
        return(i)
      }
      i <- p[i]
    }
    s[1] <- p[s[1]]
  }
  s[1]
}), sep="\n")
