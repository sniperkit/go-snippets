object CaseClass{
  case class Employee(lName:String, fName:String, ID:Int) {
    override def toString():String = {
      "Employee ID: " + ID + "\nFull Name: " + fName + " " + lName
    }
  }

  val e1 = List(Employee("Howard", "Wang", 12345), Employee("Hao", "Wang", 12346), Employee("H", "Wang", 12347))
  e1.foreach(println)
}
