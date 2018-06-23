#!/usr/bin/env python
"""
Tests for the cmakescripts.cmakemodifier module

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
import unittest
import re
import os
import glob

###
# third-party packages
# - none


###
# internal packages
import cmakemodifier
import findcmakescripts

# format for each:
# key is the filename - extension
# value is whatever the variable name suggests
inputparse = dict()
expectedoutput = dict()

def setUp():
	infiles = glob.glob(os.path.split(__file__)[0] + '/testdata/Modifications/*.parse')
	infiles.sort()
	outfiles = glob.glob(os.path.split(__file__)[0] + '/testdata/Modifications/*.output')
	outfiles.sort()

	assert len(outfiles) == len(infiles)
	assert len(infiles) == 2



	for infile, outfile in zip(infiles, outfiles):
		ibase = os.path.splitext(infile)[0]
		obase = os.path.splitext(outfile)[0]
		assert ibase == obase

		inputf = open(infile, 'r')
		inputstr = inputf.read()
		inputf.close()

		outputf = open(outfile, 'r')
		outputstr = outputf.read()
		outputf.close()

		indata = eval(inputstr)
		outdata = eval(outputstr)

	dataKeys = inputparse.keys()


## Requirement:
## A given input parse has only one expected output parse
class ExpectedCleanup(unittest.TestCase):

	subtest = ""

	if "nose" in dir():
		def _exc_info(self):
			print "Subtest info:"
			print self.subtest
			return super(ExpectedCleanup, self)._exc_info()

	def testApplyVisitor(self):
		"""passing in a known-good parse with subdirs and checking the result"""
		for key in inputparse.keys():
			inparse = inputparse[key]
			expected = expectedoutput[key]
			self.subtest = key
			self.assertEqual(cmakemodifier.apply_all_cleanup_visitors(inparse), expected)

if __name__=="__main__":
	## Run tests if executed directly
	try:
		import nose
		start = nose.main
	except (ImportError):
		start = unittest.main

	start()
