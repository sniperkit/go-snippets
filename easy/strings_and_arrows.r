cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
  v <- attr(gregexpr('(?=((>>--)>))', s, perl=T)[[1]], "match.length")
  vl <- if (identical(v, -1L)) 0 else length(v)
  w <- attr(gregexpr('(?=((<--<)<))', s, perl=T)[[1]], "match.length")
  wl <- if (identical(w, -1L)) 0 else length(w)
  vl + wl
}), sep="\n")
