object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  val cat = List("This program is for humans",
                 "Still in Mama's arms",
                 "Preschool Maniac",
                 "Elementary school",
                 "Middle school",
                 "High school",
                 "College",
                 "Working for the man",
                 "The Golden Years")
  val age = List(0, 3, 5, 12, 15, 19, 23, 66, 101)

  for (l <- lines) {
    val s = l.toInt
    println(cat((age.foldLeft(0)((a: Int, b: Int) => if (s >= b) a + 1 else a)) % 9))
  }
}
