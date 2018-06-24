object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  val pattern = "\"id\": \\d+,".r

  for (l <- lines)
    println((pattern findAllIn l).map(_.filter(_.isDigit).toInt).sum)
}
