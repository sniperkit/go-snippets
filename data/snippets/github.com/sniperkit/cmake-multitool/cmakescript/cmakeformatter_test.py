#!/usr/bin/env python
"""
Tests for the cmakescripts.cmakeformatter module

Original Author:
2010 Ryan Pavlik <rpavlik@iastate.edu> <abiryan@ryand.net>
http://academic.cleardefinition.com
Iowa State University HCI Graduate Program/VRAC
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
import cmakeparser
import cmakegrammar
import cmakeformatter
import findcmakescripts


# format for each:
# key is the filename - extension
# value is whatever the variable name suggests
inputfiles = dict()
inputstrings = dict()
inputuppers = dict()
inputlowers = dict()

parsedfiles = dict()
parsedstrings = dict()

parseduppers = dict()
parsedlowers = dict()

def setUp():
	cmakes = glob.glob(os.path.split(__file__)[0] + '/testdata/KnownValues/*.cmake')
	cmakes.sort()
	parses = glob.glob(os.path.split(__file__)[0] + '/testdata/KnownValues/*.parse')
	parses.sort()

	assert len(parses) == len(cmakes)
	assert len(parses) == 10



	for cmakefn, parsefn in zip(cmakes, parses):
		cbase = os.path.splitext(cmakefn)[0]
		pbase = os.path.splitext(parsefn)[0]
		assert cbase == pbase

		cmakef = open(cmakefn, 'r')
		cmakestr = cmakef.read()
		cmakef.close()

		parsef = open(parsefn, 'r')
		parsestr = parsef.read()
		parsef.close()

		parsedata = eval(parsestr)

		# case-change our known parses, but don't touch "None" because it's
		#
		applyExceptToNone = lambda userfunc, string, untouchable:	untouchable.join(
			[	userfunc(chunk)
				for chunk
				in string.split(untouchable)	]	)

		safeUpper = lambda input: applyExceptToNone(lambda x: x.upper(), input, "None")
		safeLower = lambda input: applyExceptToNone(lambda x: x.lower(), input, "None")

		parseupper = safeUpper(parsestr)
		parselower = safeLower(parsestr)

		inputfiles[cbase] = cmakefn
		inputstrings[cbase] = cmakestr
		inputuppers[cbase] = safeUpper(cmakestr)
		inputlowers[cbase] = safeLower(cmakestr)

		parsedfiles[cbase] = parsedata
		parsedstrings[cbase] = parsedata

		parseduppers[cbase] = parseupper
		parsedlowers[cbase] = parselower

	dataKeys = inputfiles.keys()

## Requirement:
## Given a valid parse, produce an output
class KnownParses(unittest.TestCase):

	subtest = ""

	if "nose" in dir():
		def _exc_info(self):
			print "Subtest info:"
			print self.subtest
			return super(KnownParses, self)._exc_info()

	def testCanOutputKnownParses(self):
		"""passing in a known-good parse to the formatter"""
		for key in parsedstrings.keys():
			self.subtest = key
			formatter = cmakeformatter.CMakeFormatter(parsedstrings[key])
			formatter.output_as_cmake()
			self.assertNotEqual(formatter.output_as_cmake(), "")

	def testKnownParsesRoundtrip(self):
		"""parsing formatted output should match the input parse"""
		for key in parsedstrings.keys():
			self.subtest = key
			formatter = cmakeformatter.CMakeFormatter(parsedstrings[key])
			formatted = formatter.output_as_cmake()
			self.assertEqual(cmakeparser.parse_string(formatted).parsetree, parsedstrings[key])


#class WildModules(unittest.TestCase):
#	def setUp(self):
#		basedir = os.path.split(__file__)[0] + '/testdata/WildModules'
#		self.modules = findcmakescripts.find_cmake_scripts(basedir)
#
#	def testWildModuleParsesRoundtrip(self):
#		for filename in self.modules:
#
#			try:
#				parser = cmakeparser.parse_file(filename)
#			except:
#				continue
#			formatter = cmakeformatter.CMakeFormatter(parser.parsetree)
#			formatted = formatter.output_as_cmake()
#			self.assertEqual(cmakeparser.parse_string(formatted).parsetree, parser.parsetree)


if __name__=="__main__":
	## Run tests if executed directly
	try:
		import nose
		start = nose.main
	except (ImportError):
		start = unittest.main

	start()
