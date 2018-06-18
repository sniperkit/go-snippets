#!/bin/sh

exec scala "$0" "$@"

!#

import scala.sys.process._
import scala.language.postfixOps

object Compile extends App {

  def compile () : Unit = { 
    compile("*")
  }

  def compile (expr:String) : Unit = { 
    var cmd:String = ""
    var retCode:Int = 0 

    if (expr == "*" )
      cmd = "scalac -d classes src/*.scala"
    else if (expr.endsWith(".scala"))
      cmd = "scalac -d classes " + expr
    else
      cmd = "scalac -d classes " + expr + ".scala"

    cmd += " -feature"

    print("rm -rf classes" !!) 
    print("mkdir classes" !!) 

    println("Compiling: " + cmd)
    retCode = cmd !
  }

  if (args.length == 0) {
    println("No file specified. Using *")
    compile()
  } else
    compile( args(0) )
}

Compile.main(args)
