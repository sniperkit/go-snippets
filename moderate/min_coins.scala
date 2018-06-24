object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  val b = List(0, 1, 2, 1, 2)

  for (l <- lines) {
    val n = l.toInt
    println(n / 5 + b(n % 5))
  }
}
