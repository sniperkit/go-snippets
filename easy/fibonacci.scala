object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    def nfib(c: Int): Int = {
      def nfib0(n: Int, a: Int, b: Int): Int =
        if (n == c) a else nfib0(n + 1, a + b, a)

      if (c == 0) 0 else nfib0(1, 1, 0)
    }

    println(nfib(l.toInt))
  }
}
