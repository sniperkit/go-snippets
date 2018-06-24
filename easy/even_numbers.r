cat(1 - as.numeric(readLines(tail(commandArgs(), n=1))) %% 2, sep="\n")
