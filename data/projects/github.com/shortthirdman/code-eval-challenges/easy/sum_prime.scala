object Main extends App {
  lazy val primes: Stream[Int] = 2 #:: Stream.from(3).filter(i =>
    primes.takeWhile{j => j * j <= i}.forall{k => i % k > 0})
  println(primes.take(1000).sum)
}
