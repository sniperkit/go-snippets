object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)

  for (l <- lines) {
    val s = l.split(" ").map(_.toInt)
    println(if (s(0) == s(2)) (if (s(1) == s(3)) "here" else
                               if (s(1) < s(3)) "N" else "S") else
            if (s(1) == s(3)) (if (s(0) < s(2)) "E" else "W") else
           (if (s(0) < s(2)) (if (s(1) < s(3)) "NE" else "SE") else
                             (if (s(1) < s(3)) "NW" else "SW")))
  }
}
