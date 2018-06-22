#!/usr/bin/env python
"""
Script that dumps the internal parse tree used by the CMakeScript packages.

Original Author:
2010 Ryan Pavlik <rpavlik@iastate.edu> <abiryan@ryand.net>
http://academic.cleardefinition.com
Iowa State University HCI Graduate Program/VRAC

          Copyright Iowa State University 2010.
 Distributed under the Boost Software License, Version 1.0.
    (See accompanying file LICENSE_1_0.txt or copy at
          http://www.boost.org/LICENSE_1_0.txt)
"""

###
# standard packages
import sys
import os
import subprocess
from optparse import OptionParser

###
# third-party packages
# - none

###
# internal packages
import cmakescript

class App:
	def __init__(self, args_in=sys.argv[1:]):
		self.args_in = args_in

	def main(self):
		for infile in self.args_in:
			print self.processFile(infile)

	def processFile(self, filename):
		try:
			parser = cmakescript.parse_file(filename)
		except cmakescript.IncompleteStatementError:
			print "Error parsing file: IncompleteStatementError"
			return None
		except cmakescript.UnclosedChildBlockError:
			print "Error parsing file: UnclosedChildBlockError"
			return None

		tree = cmakescript.CMakeBlock(parser.parsetree)
		return repr(tree)

###
# __main__

if __name__ == "__main__":
## Can be used as a tool when executed directly
	app = App()
	app.main()
