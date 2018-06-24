object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  val (trail, trash) = ("^[^a-z]+|[^a-z]+$".r, "[^a-z]+".r)

  for (l <- lines)
    println(trash.replaceAllIn(trail.replaceAllIn(l.toLowerCase, ""), " "))
}
