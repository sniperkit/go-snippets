object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  def dsig(a: Int): Int =
    if (a == 0) 0
    else if (a % 10 > 0) ((1 << (3 * (a % 10))) + dsig(a / 10))
    else dsig(a / 10)

  for (l <- lines) {
    val d = l.toInt
    val ds = dsig(d)
    def esig(a: Int): Int =
      if (dsig(a) == ds) a else esig(a + 9)
    println(esig(d + 9))
  }
}
