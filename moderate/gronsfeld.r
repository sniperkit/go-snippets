gro <- " !\"#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), ";"), function(s) {
  i <- 1
  paste(lapply(strsplit(s[2], "")[[1]], function(t) {
    g <- (nchar(gro) + regexpr(t, gro, fixed = TRUE)[1] - 1 - as.integer(substr(s[1], i, i))) %% nchar(gro) + 1
    i <<- i %% nchar(s[1]) + 1
    substr(gro, g, g)
  }), collapse="")
}), sep="\n")
