object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    val s = l.split(" ")
    val t = s.length - s(s.length - 1).toInt - 1
    if (t >= 0 && t < s.length)
      println(s(t))
  }
}
