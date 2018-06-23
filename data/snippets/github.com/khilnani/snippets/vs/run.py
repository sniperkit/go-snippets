#!/usr/bin/python 

import sys
import os

def run (expr):
  cmd = ""

  if expr.endswith( ".scala" ):
    expr = expr.split('.')[0]

  if expr.find ("/") != -1:
    expr = expr.split('/')[-1]

  cmd = "scala -classpath classes " + expr

  print "Running: " + cmd
  os.system (cmd)


if len(sys.argv) == 1:
  print "USAGE: " + sys.argv[0] + " CLASSNAME"
else:
  run (sys.argv[1])
