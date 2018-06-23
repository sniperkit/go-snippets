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

import glob
import os
import unittest

import parse_cmake.parsing as cmp


def yield_examples():
    paths = glob.glob(os.path.join(os.path.dirname(__file__),
                                   '..', 'example_inputs', '*'))
    for path in paths:
        with open(path) as file:
            contents = file.read()
            yield path, contents


class ExamplesTestCase(unittest.TestCase):
    def test_idempotency_of_parse_unparse(self):
        round_trip = lambda s, path='<string>': str(cmp.parse(s, path))
        for path, contents in yield_examples():
            self.assertEqual(round_trip(contents, path),
                             round_trip(round_trip(contents, path)),
                             'Failed on %s' % path)

    def test_tree_is_unchanged(self):
        for path, contents in yield_examples():
            expected = cmp.parse(contents, path)
            actual = cmp.parse(str(cmp.parse(contents, path)))
            msg = 'Failed on %s.\nExpected\n%s\n\nGot\n%s' % (path, expected, actual)
            self.assertMultiLineEqual(str(expected), str(actual), msg)
