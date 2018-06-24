object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    val s = l.split(",")
    println(if (s(0).endsWith(s(1))) 1 else 0)
  }
}
