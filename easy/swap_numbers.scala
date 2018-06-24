object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  def swp(s: String): String =
    s.takeRight(1) + s.slice(1, s.length - 1) + s.take(1)

  for (l <- lines)
    println(l.split(" ").map(swp).mkString(" "))
}
