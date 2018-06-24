cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  n <- 0
  k <- 0
  st <- sapply(strsplit(s, "")[[1]], function(t) {
    n <<- n + 1
    if (t == "$") {
      k <<- n
    }
    utf8ToInt(t) * 2048 + n
  })
  st <- st[order(st)]
  r <- c()
  for (i in 1:(n-1)) {
    r[i] <- intToUtf8(st[k] / 2048)
    k <- st[k] %% 2048
  }
  paste(r, collapse="")
}), sep="\n")
