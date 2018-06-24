object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    val s = l.split(", ").map(_.filter(_.isDigit).toInt)
    println(if (s(0) + s(1) + s(2) == 0) 0 else (s(0) * 3 + s(1) * 4 + s(2) * 5) * s(3) / (s(0) + s(1) + s(2)))
  }
}
