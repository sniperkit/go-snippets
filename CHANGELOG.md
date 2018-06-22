# 0.4.1 (2015/03/26)

- Added CI and PEP8/Flake8 Style automatted tests.
- Added testing on Python3 and fixed some issues with Python3.
- Merged pull request from @robertknight (#1), changes as he listed them:
 - Add an option to prettify() to make formatting configurable for projects that need it
 - Lower-case command names in the output for consistency
 - Change default indent to two spaces
 - Make the tests runnable with 'python -m unittest discover'
 - Refactor formatting of command arguments to make it easier to implement different formatting approaches as needed
 - Fix an issue where 'elseif()' was not indented correctly
- Clarified LICENSE as being Apache 2.0, updated source code and LICENSE file.

# 0.4.0 (2014/03/29)

- Renamed to parse_cmake from cmakelists_parsing because I cannot get a hold of the original author to get permission to release it on PyPi.
- Changed to allow for nesting left parenthesis.
