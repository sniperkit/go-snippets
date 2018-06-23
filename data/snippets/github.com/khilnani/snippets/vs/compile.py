#!/usr/bin/python

import sys
import os

def compile (expr="*"):
  print "Preparing ..." 

  os.system( "rm -rf classes" )
  os.system( "mkdir classes" )

  cmd = ''

  if expr == "*":
    cmd = "scalac -d classes src/*.scala"
  elif expr.endswith( ".scala" ):
    cmd = "scalac -d classes " + expr
  else:
    cmd = "scalac -d classes " + expr + ".scala" 

  cmd += " -feature"

  print "Compiling: " + cmd
  os.system( cmd )

  print "Done."

if len(sys.argv) < 2:
  print "No File specified. Using *"
  compile()
else:
  compile (sys.argv[1])
