object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  lazy val primes: Stream[Int] = 2 #:: Stream.from(3).filter(i =>
    primes.takeWhile{j => j * j <= i}.forall{k => i % k > 0})

  for (l <- lines)
    println(primes.takeWhile(math.pow(2, _) - 1 < l.toInt)
      .map(math.pow(2, _).toInt - 1).mkString(", "))
}
