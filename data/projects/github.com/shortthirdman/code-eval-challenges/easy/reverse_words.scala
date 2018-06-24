object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  def revw(words: List[String]): String = {
    def prevw(m: Int): String =
      if (m < 0) "" else " " + words(m) + prevw(m - 1)

    if (words.length == 0) "" else
      words(words.length - 1) + prevw(words.length - 2)
  }

  for (l <- lines)
    println(revw(l.split(" ").toList))
}
