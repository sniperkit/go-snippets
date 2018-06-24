object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  val pattern = "^((\"[0-9A-Za-z@.+-=]+\")|([0-9A-Za-z.+-=]+))@\\w*(\\w+\\.)+\\w{2,4}$".r

  for (l <- lines)
    println(l match {
      case pattern(_, _, _, _) => "true"
      case _ => "false"
    })
}
