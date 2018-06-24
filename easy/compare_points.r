cat(sapply(lapply(strsplit(readLines(tail(commandArgs(), n=1)), " ", fixed=TRUE), as.integer), function(s) {
  ifelse(s[1] == s[3], ifelse(s[2] == s[4], "here", ifelse(s[2] < s[4], "N", "S")),
         ifelse(s[2] == s[4], ifelse(s[1] < s[3], "E", "W"),
         ifelse(s[1] < s[3], ifelse(s[2] < s[4], "NE", "SE"),
                             ifelse(s[2] < s[4], "NW", "SW"))))
}), sep="\n")
