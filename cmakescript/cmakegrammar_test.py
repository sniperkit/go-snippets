#!/usr/bin/env python
"""
Tests for the cmakescript.cmakegrammar module

Original Author:
2010 Ryan Pavlik <rpavlik@iastate.edu> <abiryan@ryand.net>
http://academic.cleardefinition.com
Iowa State University HCI Graduate Program/VRAC
"""

###
# standard packages
import unittest
import re

###
# third-party packages
# - none

###
# internal packages
import cmakegrammar


## Requirement:
## All regexes must be able to be compiled
class AllRegexesCompilable(unittest.TestCase):
	def testCompileRegexes(self):
		"""Compiling regex string attributes of the grammar"""
		isReStringVar = re.compile("_re")
		reVars = filter(isReStringVar.match, dir(cmakegrammar))
		regexStrings = [vars(cmakegrammar)[varname]
						for varname in reVars
						if re.match("_re", varname)	]

		for reStr in regexStrings:
			print reStr
			re.compile(reStr)


## Requirement:
## Accept only possible command names
class AcceptRejectReFuncNames(unittest.TestCase):

	subtest = ""
	def _exc_info(self):
		print "Subtest info:"
		print self.subtest
		return unittest.TestCase._exc_info(self)

	def testAcceptValidFunctionNames(self):
		"""reFuncName: accept just valid function names"""
		data = (	"func",
				"FuNc",
				"FUNC",
				"the_func",
				"the_1_FUNC"	)
		for line in data:
			self.subtest = line
			self.assertTrue(re.search("^" + cmakegrammar._reFuncName + "$", line).group())

	def testRejectInvalidFunctionNames(self):
		"""reFuncName: reject just invalid function names"""
		data = (	"func(",
					"FuNc{",
					"FUNC-",
					"FUNC ",
					"the_func!",
					"the_1_FUNC\"",
					"#comment"	)
		for line in data:
			self.subtest = line
			result = re.search("^" + cmakegrammar._reFuncName + "$", line)
			if result is not None:
				print result.groups()
			self.assertEquals(re.search("^" + cmakegrammar._reFuncName + "$", line), None)

	def testExtractFunctionNames(self):
		"""reFuncName: extract valid function names from surroundings"""
		data = (	("func(", "func"),
					(" FuNc ", "FuNc"),
					("\tFUNC\n", "FUNC")	)
		for line, expected in data:
			self.subtest = line
			self.assertEquals(re.search(cmakegrammar._reFuncName, line).group("FuncName"), expected)

## Requirement:
## Accept only valid comments
class AcceptRejectReComment(unittest.TestCase):

	subtest = ""
	def _exc_info(self):
		print "Subtest info:"
		print self.subtest
		return unittest.TestCase._exc_info(self)

	def testAcceptValidComments(self):
		"""reComment: accept just valid comments"""
		data = (	"#",
					"#Comment",
					"# Comment",
					"## Comment",
					"###"	)
		for line in data:
			self.subtest = line
			self.assertTrue(re.search("^" + cmakegrammar._reComment + "$", line))

	def testRejectInvalidComments(self):
		"""reComment: reject just invalid comments"""
		data = (	"",
					"func(",
					"\#notacomment"	)
		for line in data:
			self.subtest = line
			self.assertEqual(re.search("^" + cmakegrammar._reComment + "$", line), None)

	def testStripCommentWS(self):
		"""reComment: strip leading and trailing whitespacefrom comments"""
		data = (	" # ",
					" # not exact comment"	)
		for line in data:
			self.subtest = line
			self.assertEqual(re.search(cmakegrammar._reComment, line).group("Comment"), line.strip())

class AcceptRejectArgument(unittest.TestCase):
	"""Parsing arguments apart using a regex"""
	def testAcceptValidArg(self):
		"""Accept a valid argument with re.match"""
		data = (	r"arg",
					r"more_args1",
					r"1arg-1"
					r"anarg\"withquote",
					r'"quoted arg"',
					r'""',
					r"''",
					r'''"this is a long
					 argument"''')
		for item in data:
			print item
			self.assertNotEqual(re.match(cmakegrammar._reArg + "$", item), None)

	def testRejectInvalidArg(self):
		"""Reject an invalid argument with re.match"""
		data = (	r"arg'",
					r" more_args1",
					r"1arg 1"
					r"anarg\"withquote",
					r'"quoted" args"',
					r'"\"')
		for item in data:
			print item
			self.assertEqual(re.match(cmakegrammar._reArg + "$", item), None)

	def testSplitArgs(self):
		"""Split valid arguments"""
		data = (	(r"arg ", 1),
					(r"more_args1", 1),
					(r"1arg 1", 2),
					(r'anarg\"withquote another', 2),
					(r'"quoted\" arg" PIE', 2),
					(r'''"this is a long
					 argument" another''', 2))
		for item, number in data:
			print item
			print re.findall(cmakegrammar._reArg, item)
			self.assertEqual(len(re.findall(cmakegrammar._reArg, item)), number)


## Requirement:
## Be able to parse a valid cmake command input line.
class ParseCompleteLine(unittest.TestCase):
	commandsOnly = (	("",				("", None, None)),
					("func()",			("func", None, None)),
					("func(arg)",		("func", "arg", None)),
					("func(arg arg)",	("func", "arg arg", None)),
					("func( arg arg )",	("func", "arg arg", None)),
					("func( arg  ar )",	("func", "arg  ar", None)),
					(r"func(\#notcmnt)",("func", r"\#notcmnt", None)),
					("func(oh\nyes\nmultiline)", ("func", "oh yes multiline", None))	)

	commentsOnly = (	"#",
				"# comment",
				"#comment",
				"## comment",
				"##comment",
				"#a#a#a#"	)

	mixed = (	("func() # cmnt",	("func", None, "# cmnt")),
			(r"func(\#notcmnt) #iscmnt",("func", r"\#notcmnt", "#iscmnt")),
			("func(arg1 #comment\narg2)", ("func", "arg1 arg2", "#comment")),
			("func(arg1 # comment morecomment\narg2)", ("func", "arg1 arg2", "# comment morecomment")),
			("func(arg1 # comment morecomment\narg2) # more comment", ("func", "arg1 arg2", "# comment morecomment\n# more comment")))



	subtest = ""
	def _exc_info(self):
		print "Subtest info:"
		print self.subtest
		return unittest.TestCase._exc_info(self)

	def testParseLineCommand(self):
		"""parse_line on a line with a command only"""
		for (line, expected) in self.commandsOnly:
			self.subtest = line
			func, args, comment = cmakegrammar.parse_line(line)


			self.assertEqual((func, args, comment), expected)

	def testParseLineComment(self):
		"""parse_line on a line with a comment only"""
		for commentstring in self.commentsOnly:
			self.subtest = commentstring
			func, args, comment = cmakegrammar.parse_line(commentstring)

			self.assertEqual((func, args, comment), ("", None, commentstring))

	def testParseLineMixed(self):
		"""parse_line on a line with both a command and a comment"""
		for (line, expected) in self.mixed:
			self.subtest = line
			func, args, comment = cmakegrammar.parse_line(line)

			self.assertEqual((func, args, comment), expected)

## Requirement:
## Notify on an incomplete cmake input line
class ParsePartialLine(unittest.TestCase):
	data = (	"func(",
				"func(arg",
				"func(arg ",
				"func ( arg ",
				"func(arg arg"	)


	subtest = ""
	def _exc_info(self):
		print "Subtest info:"
		print self.subtest
		return unittest.TestCase._exc_info(self)

	def testParseIncompleteStatement(self):
		"""parse_line on an incomplete line with a command only"""
		for line in self.data:
			self.subtest = line
			self.assertRaises(cmakegrammar.IncompleteStatementError,
							cmakegrammar.parse_line, line)

class HandleEOFSentry(unittest.TestCase):
	def testHandleEOFSentry(self):
		"""Parsing None as your line results in an all-None tuple"""
		self.assertEqual(cmakegrammar.parse_line(None), (None, None, None))


if __name__=="__main__":
	## Run tests if executed directly
	try:
		import nose
		start = nose.main
	except (ImportError):
		start = unittest.main

	start()
