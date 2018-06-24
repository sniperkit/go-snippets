robot <- function(f, x, y) {
  if (x == 3 && y == 3) {
    return(1)
  }
  ret <- 0
  if (x > 0 && !((x-1 + 4*y) %in% f)) {
    ret = ret + robot(append(f, x-1 + 4*y), x-1, y)
  }
  if (y > 0 && !((x + 4*(y-1)) %in% f)) {
    ret = ret + robot(append(f, x + 4*(y-1)), x, y-1)
  }
  if (x < 3 && !((x+1 + 4*y) %in% f)) {
    ret = ret + robot(append(f, x+1 + 4*y), x+1, y)
  }
  if (y < 3 && !((x + 4*(y+1)) %in% f)) {
    ret = ret + robot(append(f, x + 4*(y+1)), x, y+1)
  }
  ret
}

cat(robot(c(0), 0, 0), sep="\n")
