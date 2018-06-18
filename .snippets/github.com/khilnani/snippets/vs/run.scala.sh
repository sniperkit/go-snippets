#!/bin/sh

exec scala "$0" "$@"

!#

import scala.sys.process._
import scala.language.postfixOps

object Run extends App {

  def run (expr:String) : Unit = { 
    var cmd:String = ""
    // function parameters are immutable in scala, define a local var
    var cname:String = expr 
    var retCode:Int = 0 

    if (cname.endsWith(".scala"))
      cname = cname.split('.')(0)
    
    if (cname.indexOf('/') > -1)
      cname = cname.split('/').last

    cmd = "scala -classpath classes " + cname

    println("Running: " + cmd)
    retCode = cmd !
  }

  if (args.length == 0)
    println("USAGE: run.scala.sh CLASSNAME")
  else
    run( args(0) )
}

Run.main(args)
