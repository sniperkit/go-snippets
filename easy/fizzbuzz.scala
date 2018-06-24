object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  def fzbz(f: Int, b: Int, n: Int): String = {
    def fb(m: Int): String = {
      if (m > n) "" else " " + {
        if (m % f == 0 && m % b == 0) "FB" + fb(m + 1) else {
          if (m % f == 0) "F" + fb(m + 1) else {
            if (m % b == 0) "B" + fb(m + 1) else
              m.toString + fb(m + 1)
          }
        }
      }
    }

    if (n < 1) "" else {
      if (f == 1 && b == 1) "FB" + fb(2) else {
        if (f == 1) "F" + fb(2) else {
          if (b == 1) "B" + fb(2) else
            "1" + fb(2)
        }
      }
    }
  }

  for (l <- lines) {
    val fbn = l.split(" ").map(_.toInt)
    println(fzbz(fbn(0), fbn(1), fbn(2)))
  }
}
