cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  paste(lapply(strsplit(s, "")[[1]], function(x) {
    if (x %in% letters[1:6]) {
      rawToChar(as.raw(strtoi(charToRaw(x), 16L) + 20))
    } else if (x %in% letters[7:13]) {
      rawToChar(as.raw(strtoi(charToRaw(x), 16L) + 7))
    } else if (x %in% letters[14:20]) {
      rawToChar(as.raw(strtoi(charToRaw(x), 16L) - 7))
    } else if (x %in% letters[21:26]) {
      rawToChar(as.raw(strtoi(charToRaw(x), 16L) - 20))
    } else {
      x
    }
  }), collapse="")
}), sep="\n")
