#!/usr/bin/env python
"""
Recursively find CMakeLists.txt and *.cmake

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
import os.path
import re

###
# third-party packages
# - none

###
# internal packages
from recursivelistdir import recursive_listdir_iter

def find_cmake_scripts(startPath):

	# Get the path the way we want it.
	startPath = os.path.abspath(startPath)

	if os.path.isfile(startPath):
		# Any file directly passed in is assumed to be a script, no matter its name
		return [startPath]


	# Build a regex to find CMakeLists.txt or *.cmake, case insensitive
	# A match against the filename (not including the path!) means that
	# this is probably a valid input file.
	reCMakeLists = r"(?ix)^(CMakeLists\.txt)$"
	reCMakeModule = r"(?ix)(\.cmake)$"
	isScript = re.compile(r"(" + reCMakeLists + r"|" + reCMakeModule + r")")

	# A little function to use with filter - so we only apply the regex
	# to the file's basename, not the full path
	isPathScript = lambda filepath: isScript.search(os.path.basename(filepath))

	# Now, actually do it.
	# Use our regex to filter the recursive list of files returned
	cmakeScripts = filter(isPathScript, recursive_listdir_iter(startPath))

	return cmakeScripts
