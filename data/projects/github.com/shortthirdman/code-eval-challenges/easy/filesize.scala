import java.io.File

object Main extends App {
  val f = new File(args(0))
  println(f.length)
}
