object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    def binomial(n: Int, k: Int) =
      (BigInt(n - k + 1) to n).product / (BigInt(1) to k).product
    val m = l.toInt
    var r = BigInt(0)
    for (i <- 0 to m/2)
      r += binomial(m, i) * binomial(11*m/2-1-10*i, m-1) * (scala.math.pow(-1, i)).toInt
    println(r)
  }
}
