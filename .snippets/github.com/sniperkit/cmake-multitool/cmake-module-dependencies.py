#!/usr/bin/env python
"""
Main application to use the CMakeScript packages to output cleaner code.

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
import yapgvb	# GraphViz

###
# internal packages
import cmakescript
from mergetool import MergeTool


class App:
	def __init__(self, args_in=sys.argv[1:]):
		self.args_in = args_in
		self.mergetool = None

	def main(self):
		parser = OptionParser(usage="usage: %prog [options] [[file|dir]...]",
							  version="%prog 0.5, part of the cmakescript tools")

		parser.add_option("-m", "--merge",
							type="choice",
							choices=MergeTool.mergetools.keys(),
							metavar="APPNAME",
							dest="mergetool",
							default=None,
							help="open a diff/merge app APPNAME for each file "
								 "processed.  Supported APPNAME options are: " +
								 " ".join(MergeTool.mergetools.keys())
							)

		parser.add_option("-q", "--quiet",
						action="store_false", dest="verbose", default=True,
						help="don't print status messages to stdout")

		(self.options, args) = parser.parse_args(self.args_in)

		if len(args) == 0:
			args.append(os.getcwd())

		inputfiles = []

		for arg in args:
			inputfiles.extend(cmakescript.find_cmake_scripts(arg))

		dependencies["findmodules"] = {}
		dependencies["othermodules"] = {}
		dependencies["optionalmodules"] = {}
		dependencies["files"] = {}
		dependencies["optionalfiles"] = {}

		knownfiles = [os.path.relpath(x) for x in inputfiles]
		justname = lambda x:os.path.splitext(os.path.basename(x))[0]
		allmodules["found"] = [justname(x) for x in knownfiles]
		t1 = subprocess.Popen(["cmake", "--help-modules-list"], stdout=subprocess.PIPE)
		allmodules["system"] = [x.strip()
								for x in t1.communicate()[0].splitlines()[1:]]

		nodes = {}
		edges = []

		for key, val in allmodules.iteritems():
			findmodules[key] = [x
								for x in allmodules[key]
								if re.match(r"Find", x)]
			othermodules[key] = [x
								 for x in allmodules[key]
								 if re.match(r"Find", x) is None]

		for infile, number in zip(inputfiles, range(1, len(inputfiles)+1)):
			print "------------------------"
			print infile + " - " + str(number) + " of " + str(len(inputfiles))
			print "------------------------"

			visitor = self.processFile(infile)
			shortname = os.path.relpath(infile)
			pathto = os.path.split(shortname)[0]
			

			dependencies["findmodules"][shortname] = visitor.findmodules
			dependencies["othermodules"][shortname] = visitor.modules
			dependencies["othermodules"][shortname] = visitor.modules
			dependencies["optionalmodules"][shortname] = visitor.optionalmodules
			dependencies["files"][shortname] = [os.path.join(pathto, x)
												for x in visitor.files]
			dependencies["optionalfiles"][shortname] = [os.path.join(pathto, x)
														for x in visitor.optionalfiles]




	def processFile(self, filename):
		try:
			parser = cmakescript.parse_file(filename)
		except cmakescript.IncompleteStatementError:
			print "Error parsing file: IncompleteStatementError"
			return None
		except cmakescript.UnclosedChildBlockError:
			print "Error parsing file: UnclosedChildBlockError"
			return None

		visitor = cmakescript.VisitorFindModuleDependencies()
		tree = cmakescript.CMakeBlock(parser.parsetree)
		tree.accept(visitor)
		return visitor


###
# __main__

if __name__ == "__main__":
## Can be used as a tool when executed directly
	app = App()
	app.main()
