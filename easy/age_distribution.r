categories <- c("This program is for humans", "Still in Mama's arms", "Preschool Maniac", "Elementary school", "Middle school", "High school", "College", "Working for the man", "The Golden Years")
ages <- c(0, 3, 5, 12, 15, 19, 23, 66, 101)

cat(sapply(as.integer(readLines(tail(commandArgs(), n=1))), function(s) {
  i <- 1
  while (i <= length(ages) && s >= ages[i]) {
    i <- i + 1
  }
  if (i > length(ages)) {
    i <- 1
  }
  categories[i]
}), sep="\n")
