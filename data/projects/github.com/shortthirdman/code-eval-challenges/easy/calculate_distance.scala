object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    val t = l.filter(" -0123456789".contains(_)).split(" ").map(_.toInt)
    val (x, y) = (t(0) - t(2), t(1) - t(3))
    println(math.sqrt(x*x + y*y).toInt)
  }
}
