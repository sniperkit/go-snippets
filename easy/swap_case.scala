object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  def swp(c: Char): Char =
    if (c.isLower) c.toUpper else c.toLower

  for (l <- lines)
    println(l.map(swp))
}
