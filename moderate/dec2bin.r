cat(Filter(function(x) !is.na(x), sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  if (s == "") {
    return(NA)
  } else {
    gsub("^0+", "0", gsub("^0+1", "1", paste(rev(as.integer(intToBits(s))), collapse="")))
  }
})), sep="\n")
