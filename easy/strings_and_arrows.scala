object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  val (p1, p2) = ("(?=>>-->)".r, "(?=<--<<)".r)

  for (l <- lines)
    println(p1.findAllIn(l).length + p2.findAllIn(l).length)
}
