object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  val dig = Map("zero" -> 0, "one" -> 1, "two" -> 2, "three" -> 3, "four" -> 4,
                "five" -> 5, "six" -> 6, "seven" -> 7, "eight" -> 8, "nine" -> 9)

  for (l <- lines)
    println(l.split(";").map(dig).mkString)
}
