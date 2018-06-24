cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  balanced <- function(x = 1, y = 0) {
    if (y < 0) {
      return(FALSE)
    }
    i <- regexpr("[:()]", substring(s, x))[1] - 1
    if (i < 0) {
      return(y == 0)
    }
    x <- x + i
    if (substr(s, x, x) == "(") {
      balanced(x+1, y+1)
    } else if (substr(s, x, x) == ")") {
      balanced(x+1, y-1)
    } else if (substr(s, x, x) == ":") {
      if (x < nchar(s) & substr(s, x+1, x+1) == "(") {
        balanced(x+2, y) | balanced(x+2, y+1)
      } else if (x < nchar(s) & substr(s, x+1, x+1) == ")") {
        balanced(x+2, y) | balanced(x+2, y-1)
      } else {
        balanced(x+1, y)
      }
    } else {
      FALSE
    }
  }
  if (balanced()) {
    "YES"
  } else {
    "NO"
  }
}), sep="\n")
