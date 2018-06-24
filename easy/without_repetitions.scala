object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines)
    println(l.foldLeft(l(0).toString)((x: String, y: Char) => if (x.last == y) x else x + y))
}
