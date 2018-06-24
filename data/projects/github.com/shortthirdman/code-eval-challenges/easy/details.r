for (line in readLines(tail(commandArgs(), n=1))) {
  cat(sapply(strsplit(line, ","), function(s) {
    min(attr(regexpr("X\\.*Y", s), "match.length")) - 2
  }), sep="\n")
}
