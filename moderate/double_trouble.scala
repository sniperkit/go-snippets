import scala.util.control._

object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  val loop = new Breaks

  for (l <- lines) {
    val n = l.length / 2
    var r = 1
    loop.breakable {
      for (i <- 0 to n-1) {
        if ((l(i) == 'A' && l(i+n) == 'B') || (l(i) == 'B' && l(i+n) == 'A')) {
          r = 0
          loop.break
        } else if (l(i) == '*' && l(i+n) == '*') {
          r *= 2
        }
      }
    }
    println(r)
  }
}
