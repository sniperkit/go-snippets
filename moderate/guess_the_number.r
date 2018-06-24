cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(s) {
  lo <- 0
  hi <- as.integer(s[1])
  for (i in s[2:length(s)]) {
    if (i == "Lower") {
      hi <- ceiling((lo+hi)/2)-1
    } else if (i == "Higher") {
      lo <- ceiling((lo+hi)/2)+1
    }
  }
  ceiling((lo+hi)/2)
}), sep="\n")
