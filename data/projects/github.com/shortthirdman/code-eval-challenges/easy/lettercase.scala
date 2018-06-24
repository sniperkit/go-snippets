object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    val lo = l.count(_.isLower).toFloat
    val up = l.count(_.isUpper).toFloat
    if (lo + up > 0)
      print("lowercase: %.2f uppercase: %.2f\n".format(100*lo/(lo+up), 100*up/(lo+up)))
    else
      println("lowercase: 0.00 uppercase: 0.00")
  }
}
