cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  a <- paste(lapply(strsplit(s, "")[[1]], function(x) {
    if (x %in% c('0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0')) {
      x
    } else if (x %in% letters[1:10]) {
      toString(utf8ToInt(x) - 97)
    } else {
      ""
    }
  }), collapse="")
  if (a == "") {
    "NONE"
  } else {
    a
  }
}), sep="\n")
