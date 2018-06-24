object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  val m = List(50, 25, 10, 5, 1)

  def alter(n: Int, p: Int): Int = {
    if (n == 0) 1 else {
      var q = p
      while (m(q) > n) {
        q += 1
      }
      if (m(q) == 1) 1 else alter(n - m(q), q) + alter(n, q + 1)
    }
  }

  for (l <- lines)
    println(alter(l.toInt, 0))
}
