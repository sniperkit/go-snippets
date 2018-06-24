phrases <- c(", yeah!", ", this is crazy, I tell ya.", ", can U believe this?", ", eh?", ", aw yea.", ", yo.", "? No way!", ". Awesome!")
o <- 1
f <- FALSE

cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  t <- gregexpr(pattern='[.!?]', s)[[1]]
  l <- 1
  u <- ""
  for (i in 1:length(t)) {
    if (t[i] > 1) {
      u <- paste(u, substr(s, l, t[i] - 1), sep="")
    }
    if (f) {
      u <- paste(u, phrases[o], sep="")
      o <<- o %% 8 + 1
    } else {
      u <- paste(u, substr(s, t[i], t[i]), sep="")
    }
    l <- t[i] + 1
    f <<- !f
  }
  if (l < nchar(s)) {
    u <- paste(u, substr(s, l, nchar(s)), sep="")
  }
  u
}), sep="\n")
