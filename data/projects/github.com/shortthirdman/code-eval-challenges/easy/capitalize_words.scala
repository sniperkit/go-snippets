object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    val s = l.split(" ")
    val r = for (i <- 0 to (s.length - 1)) yield s(i)(0).toUpper + s(i).substring(1, s(i).length)
    println(r.mkString(" "))
  }
}
