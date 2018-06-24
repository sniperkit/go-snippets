cat(sapply(readLines(tail(commandArgs(), n=1)), function(s) {
	if (nchar(s) > 55) {
		s <- substr(s, 1, 40)
		u <- regexpr(" [^ ]*$", s)
		if (u > 0) {
			s <- substr(s, 1, u-1)
		}
		s <- paste(s, "... <Read More>", sep="")
	}
	s
}), sep="\n")
