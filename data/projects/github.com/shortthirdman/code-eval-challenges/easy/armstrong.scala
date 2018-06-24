object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  def isarmst(s: Int): String = {
    val e = math.log10(s).toInt + 1
    var q, r = s
    while (q > 0) {
      r = r - math.pow(q % 10, e).toInt
      q = q / 10
    }
    if (r == 0) "True" else "False"
  }

  for (l <- lines)
    println(isarmst(l.toInt))
}
