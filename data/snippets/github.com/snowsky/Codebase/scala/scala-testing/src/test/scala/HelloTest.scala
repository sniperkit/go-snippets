package com.scala.test

import org.scalatest.FunSuite
import main.scala._

class HelloTests extends FunSuite {
  test("the name is set correctly in constructor") {
    val p = Person("Barney Rubble")
    assert(p.name == "Barney Rubble")
  }

  test("a Person's name can be changed") {
    val p = Person("Chad Johnson")
    p.name = "Ochocinco"
    assert(p.name == "Ochocinco")
  }

  test("github api is working") {
    var p = GAPI()
    assert(p.getStatus() == true)
  }
}
