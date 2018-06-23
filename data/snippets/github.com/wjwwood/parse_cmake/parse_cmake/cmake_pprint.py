# Copyright 2015 Open Source Robotics Foundation, Inc.
# Copyright 2013 Willow Garage, Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

from __future__ import print_function

import argparse
import sys

from .parsing import parse


def main():
    # Parse arguments
    parser = argparse.ArgumentParser(description='Pretty-print CMakeLists files.')
    parser.add_argument('files', type=str, nargs='*',
                        help='files to pretty print (default is stdin)')
    parser.add_argument('-t', '--tree', action='store_true',
                        help='print out the syntax trees')
    args = parser.parse_args()

    # Gather files
    filenames = args.files
    files = [('<stdin>', sys.stdin)]
    if filenames:
        files = [(name, open(name)) for name in filenames]

    # Process files
    for (name, file) in files:
        with file:
            input = file.read()
            tree = parse(input, path=name)
            if args.tree:
                # Print out AST
                print(repr(tree))
            else:
                # Pretty print
                print(str(tree), end='')

if __name__ == '__main__':
    main()
