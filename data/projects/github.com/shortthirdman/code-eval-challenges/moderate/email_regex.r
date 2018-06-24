email <- "^((\"[0-9A-Za-z@.+-=]+\")|([0-9A-Za-z.+-=]+))@[0-9A-Za-z_]*([0-9A-Za-z_]+[.])+[0-9A-Za-z_]{2,4}$"

cat(sapply(grepl(email, readLines(tail(commandArgs(), n=1))), function(b) {
  if (b) return("true")
  "false"
}), sep="\n")
