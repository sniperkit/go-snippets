object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    def nfib(c: Int): BigInt = {
      def nfib0(n: Int, a: BigInt, b: BigInt): BigInt =
        if (n == c) a else nfib0(n + 1, a + b, a)

      if (c < 4) c else nfib0(1, 1, 1)
    }

    println(nfib(l.toInt))
  }
}
