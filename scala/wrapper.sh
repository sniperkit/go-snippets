#!/bin/sh
exec scala -save $0 $@
::!#

import java.io.File

object App {

  def main(args: Array[String]): Unit = {
      println("Hello, " + args(0))
        }
}
