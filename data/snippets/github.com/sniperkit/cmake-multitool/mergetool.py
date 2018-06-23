#!/usr/bin/env python
"""
Abstract calling of merge tools cross-platform.

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
import subprocess

###
# third-party packages
# - none

###
# internal packages
# - none

class MergeTool:
	mergetools = {	"diffmergemac" :	[
						"/Applications/DiffMerge.app/Contents/MacOS/DiffMerge",
						"-t1='Decrufted'",
						"-t2='Result'",
						"-t3='Original'",
						"{L}",
						"{C}",
						"{R}"],
					"meld"	:	[
						"meld",
						"--diff",
						"{L}",
						"{C}",
						"{R}"],
					"diffuse"	:	[
						"diffuse",
						"{L}",
						"{C}",
						"{R}"]}

	def __init__(self, tool):
		self.tool = self.mergetools[tool]

	def run(self, left, center, right):
		args = [	arg.format(L=left, R=right, C=center)
					for arg in self.tool	]
		return subprocess.call(args)

###
# __main__

#if __name__ == "__main__":
#	pass
