cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  l <- strtoi(charToRaw(substr(s, 1, 1)), 16L) - 97
  r <- strtoi(charToRaw(substr(s, 2, 2)), 16L) - 49
  m <- c()
  if (l > 1) {
    if (r > 0) {
      m <- append(m, paste(rawToChar(as.raw(l-2+97)), rawToChar(as.raw(r-1+49)), sep=""))
    }
    if (r < 7) {
      m <- append(m, paste(rawToChar(as.raw(l-2+97)), rawToChar(as.raw(r+1+49)), sep=""))
    }
  }
  if (l > 0) {
    if (r > 1) {
      m <- append(m, paste(rawToChar(as.raw(l-1+97)), rawToChar(as.raw(r-2+49)), sep=""))
    }
    if (r < 6) {
      m <- append(m, paste(rawToChar(as.raw(l-1+97)), rawToChar(as.raw(r+2+49)), sep=""))
    }
  }
  if (l < 7) {
    if (r > 1) {
      m <- append(m, paste(rawToChar(as.raw(l+1+97)), rawToChar(as.raw(r-2+49)), sep=""))
    }
    if (r < 6) {
      m <- append(m, paste(rawToChar(as.raw(l+1+97)), rawToChar(as.raw(r+2+49)), sep=""))
    }
  }
  if (l < 6) {
    if (r > 0) {
      m <- append(m, paste(rawToChar(as.raw(l+2+97)), rawToChar(as.raw(r-1+49)), sep=""))
    }
    if (r < 7) {
      m <- append(m, paste(rawToChar(as.raw(l+2+97)), rawToChar(as.raw(r+1+49)), sep=""))
    }
  }
  paste(m, collapse=" ")
}), sep="\n")
