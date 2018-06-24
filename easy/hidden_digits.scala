object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    val t = for (s <- l.filter("0123456789abcdefghij".contains(_))) yield if(s > '9') (s.toInt - 49).toChar else s
    println(if(t.isEmpty) "NONE" else t)
  }
}
