cat(gsub("([0-9]+)([^ 0-9]+)([0-9]+)", "\\3\\2\\1", readLines(tail(commandArgs(), n=1))), sep="\n")
