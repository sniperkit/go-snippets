for (i in 1:12) {
  for (j in 1:12) {
    cat(format(i*j, width=4))
  }
  cat('\n')
}
