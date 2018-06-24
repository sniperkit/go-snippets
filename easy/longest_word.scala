object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines)
    println(l.split(" ").foldLeft("")((x: String, y: String) => if (y.length > x.length) y else x))
}
