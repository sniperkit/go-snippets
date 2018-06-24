morse <- "ETIANMSURWDKGOHVF L PJBXCYZQ  54 3   2       16       7   8 90"

cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  paste(lapply(strsplit(s, " ")[[1]], function(x) {
    if (x == "") {
      " "
    } else {
      p <- Reduce(function(a, b) {
        if (b == '.') {
          a * 2
        } else {
          a * 2 + 1
        }
      }, strsplit(x, "")[[1]], 1)
      substr(morse, p-1, p-1)
    }
  }), collapse="")
}), sep="\n")
