cat(sapply(tolower(readLines(tail(commandArgs(), n=1))), function(s) {
  t <- letters[!letters %in% strsplit(s, "")[[1]]]
  if(length(t) == 0) "NULL" else paste(t, collapse="")
}), sep="\n")
