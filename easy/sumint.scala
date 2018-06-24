object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  println(lines.foldLeft(0)((a: Int, b: String) => a + b.toInt))
}
