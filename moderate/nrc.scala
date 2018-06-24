object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    val s = for (i <- l if (l.count(_ == i) == 1)) yield i
    println(if(s.isEmpty) "" else s(0))
  }
}
