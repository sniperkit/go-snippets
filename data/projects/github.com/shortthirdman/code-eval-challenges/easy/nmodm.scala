object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    val p = l.split(",").map(_.toInt)
    println(p(0) - (p(0)/p(1)) * p(1))
  }
}
