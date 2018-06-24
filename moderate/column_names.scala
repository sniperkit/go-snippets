object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  def colu(s: Int): String =
    if (s == 0) "" else colu((s - 1) / 26) ++ ('A' + (s - 1) % 26).toChar.toString

  for (l <- lines)
    println(colu(l.toInt))
}
