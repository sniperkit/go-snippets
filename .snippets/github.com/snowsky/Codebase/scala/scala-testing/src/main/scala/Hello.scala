package main.scala

object Hello extends App {
  val p = Person("Avlin Alexander")
  println("Hello from " + p.name)
}

case class Person(var name: String)
case class GAPI {
  def getStatus(): Boolean = true
}
