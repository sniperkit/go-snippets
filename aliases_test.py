import unittest

from os import path
import aliases
import json

class TestAliases(unittest.TestCase):
    def test_aliases(self):
        testdir = path.join(path.dirname(path.abspath(__file__)), 'testdata')
        rootdir = path.join(testdir, 'reexporter')
        a = aliases.get_aliases(rootdir)

        actual = json.dumps(a, sort_keys=True, indent=2, separators=': ')
        with open(path.join(rootdir, 'test_aliases.out'), 'w') as f:
            f.write(actual)
        with open(path.join(rootdir, 'test_aliases.exp')) as f:
            expected = f.read()
        self.assertEquals(actual, expected)
