object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  val n = lines.next().toInt
  println(lines.toList.sortWith(_.length > _.length).take(n).mkString("\n"))
}
