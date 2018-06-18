import sys.process._

while(true) {
  print("> ")
  var cmd = scala.io.StdIn.readLine()

  if(cmd.length > 0) {
    try {
      cmd !
    } catch {
      //case ex:
      case x: Throwable => println("Wrong shell command!")
    } finally {
    }
  }
}
