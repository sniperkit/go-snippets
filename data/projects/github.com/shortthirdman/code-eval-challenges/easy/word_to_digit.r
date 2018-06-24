w <- c(0:9)
names(w) <- c("zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine")

cat(sapply(strsplit(readLines(tail(commandArgs(), n=1)), ";"), function(s) {
  paste(lapply(s, function(x) toString(w[x])), collapse="")
}), sep="\n")
