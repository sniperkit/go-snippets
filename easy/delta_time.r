cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), " "), function(s) {
  d <- as.integer(strsplit(s[1], ":")[[1]]) - as.integer(strsplit(s[2], ":")[[1]])
  t <- abs(d[1]*3600 + d[2]*60 + d[3])
  sprintf("%02d:%02d:%02d", as.integer(t/3600), as.integer(t/60) %% 60, t %% 60)
}), sep="\n")
