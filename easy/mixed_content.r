cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), ","), function(s) {
  r = ""
  nums <- c()
  strs <- c()
  for (i in s) {
    if (is.na(as.integer(i))) {
      strs <- append(strs, i)
    } else {
      nums <- append(nums, i)
    }
  }
  if (length(strs) > 0) {
    r = paste(strs, collapse=",")
    if (length(nums) > 0) {
      r = paste(r, "|", sep="")
    }
  }
  if (length(nums) > 0) {
    r = paste(r, paste(nums, collapse=","), sep="")
  }
  r
}), sep="\n")
