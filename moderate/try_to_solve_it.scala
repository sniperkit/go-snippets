object Main extends App {
  val source = scala.io.Source.fromFile(args(0))
  val lines = source.getLines.filter(_.length > 0)
  def swp(c: Char): Char = {
    c match {
      case 'a'|'b'|'c'|'d'|'e'|'f'     => (c + 20).toChar
      case 'g'|'h'|'i'|'j'|'k'|'l'|'m' => (c + 7).toChar
      case 'n'|'o'|'p'|'q'|'r'|'s'|'t' => (c - 7).toChar
      case 'u'|'v'|'w'|'x'|'y'|'z'     => (c - 20).toChar
      case _                           => c
    }
  }

  for (l <- lines)
    println(l.map(swp))
}
