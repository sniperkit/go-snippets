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
# - none

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

		for infile, number in zip(inputfiles, range(1, len(inputfiles)+1)):
			print "------------------------"
			print infile + " - " + str(number) + " of " + str(len(inputfiles))
			print "------------------------"

			output = self.processFile(infile)
			if output is not None:
				# A trailing newline
				output = output + "\n"
				self.runMergeTool(infile, output)
				#if self.mergetool is not None and len(inputfiles) > 1:

				#	x = raw_input("Press enter to continue to the next file")


	def processFile(self, filename):
		try:
			parser = cmakescript.parse_file(filename)
		except cmakescript.IncompleteStatementError:
			print "Error parsing file: IncompleteStatementError"
			return None
		except cmakescript.UnclosedChildBlockError:
			print "Error parsing file: UnclosedChildBlockError"
			return None

		#formatter = cmakescript.CMakeFormatter(parser.parsetree)
		cleaned = cmakescript.apply_all_cleanup_visitors(parser.parsetree)
		formatter = cmakescript.NiceFormatter(cleaned)
		return formatter.output_as_cmake()

	def runMergeTool(self, filename, formatted):
		if self.mergetool is None and self.options.mergetool is not None:
			self.mergetool = MergeTool(self.options.mergetool)

		orig = open(filename, 'r')
		originalscript = orig.read()
		orig.close()
		if originalscript == formatted:
				return

		if self.mergetool is not None:

			modname = os.path.splitext(os.path.basename(filename))[0]

			t1 = subprocess.Popen(["mktemp", "-d"], stdout=subprocess.PIPE)
			tempdir = t1.communicate()[0].strip()
			tempclean = os.path.join(tempdir, modname+".Decrufted.cmake")
			temporig = os.path.join(tempdir, modname+".Original.cmake")


			temporigfile = open(temporig, 'w', False)
			temporigfile.write(originalscript)
			temporigfile.close()

			tempcleanfile = open(tempclean, 'w', False)
			tempcleanfile.write(formatted)
			tempcleanfile.close()

			self.mergetool.run(tempclean, filename, temporig)
		else:
			# If we aren't merging, print the formatted output
			print formatted

###
# __main__

if __name__ == "__main__":
## Can be used as a tool when executed directly
	app = App()
	app.main()
