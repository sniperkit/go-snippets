object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  def happy(a: Int): Int = {
    def divten(b: Int): Int = {
      if (b < 10) 0 else b / 10
    }

    if (a < 10) a * a else
      (a % 10) * (a % 10) + happy(divten(a))
  }

  def ishappy(a: Int, l: List[Int]): Int = {
    if (a == 1) 1 else {
      if (l contains a) 0 else
        ishappy(happy(a), a :: l)
    }
  }

  for (l <- lines)
    println(ishappy(l.toInt, List()))
}
