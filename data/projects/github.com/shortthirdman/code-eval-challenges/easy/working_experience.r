m <- c(1:12)
names(m) <- c("Jan", "Feb", "Mar", "Apr", "May", "Jun",
              "Jul", "Aug", "Sep", "Oct", "Nov", "Dec")
month <- function(s) {
  12 * (as.numeric(substr(s, 5, 8)) - 1990) + m[substr(s, 1, 3)] - 1
}

cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), "; "), function(s) {
  r <- vector()
  for (i in s) {
    for (j in month(substr(i, 1, 8)):month(substr(i, 10, 17))) {
      if (!(j %in% r)) {
        r <- append(r, j)
      }
    }
  }
  length(r) %/% 12
}), sep="\n")
