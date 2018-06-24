import scala.util.matching.Regex

object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  val pattern = new Regex("X[.]*Y")

  for (l <- lines)
    println((pattern findAllIn l).map(_.length).min - 2)
}
