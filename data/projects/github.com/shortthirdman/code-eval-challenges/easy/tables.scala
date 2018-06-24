object Main extends App {
  var i = 0;
  var j = 0;
  for (i <- 1 to 12; j <- 1 to 12) {
    print("%4d".format(i*j));
    if (j == 12) { println; }
  }
}
