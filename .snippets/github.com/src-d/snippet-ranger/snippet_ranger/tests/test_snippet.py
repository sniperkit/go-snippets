import unittest

import asdf
import numpy as np
from ast2vec.bblfsh_roles import Node

from modelforge import split_strings
from snippet_ranger.models.snippet import Snippet
from snippet_ranger.tests import models


def validate_asdf_file(obj, filename):
    data = asdf.open(filename)
    obj.assertIn("meta", data.tree)
    obj.assertIn("filenames", data.tree)
    obj.assertIn("sources", data.tree)
    obj.assertIn("uasts", data.tree)
    obj.assertIn("repository", data.tree)
    obj.assertIn("positions", data.tree)
    Node.FromString(split_strings(data.tree["uasts"])[0])
    obj.assertEqual(data.tree["sources"]["lengths"].shape[0],
                    data.tree["uasts"]["lengths"].shape[0])
    obj.assertEqual(0, len(data.tree["meta"]["dependencies"]))
    obj.assertEqual(data.tree["meta"]["model"], "snippet")


class TestSnippet(unittest.TestCase):
    def setUp(self):
        self.model = Snippet().load(source=models.SNIPPET)

    def assert_np_arrays(self, arr1, arr2):
        self.assertEqual(type(arr1), np.ndarray)
        self.assertEqual(type(arr2), np.ndarray)
        self.assertTrue(np.all(arr1 == arr2))

    def test_names(self):
        self.assertEqual(self.model.names,
                         ["./snippet_ranger/snippet_ranger/tests/data/test_repo/example.py_4_11"])

    def test_item(self):
        model_item = self.model[0]
        self.assertEqual(len(model_item), 4)
        self.assertEqual(model_item[0], "example.py")
        self.assertEqual(type(model_item[1]), Node)
        self.assertEqual(model_item[2],
                         "def f():\n    f2()\n    f2()\n    f3()\n    f1()\n    f3()\n    f3()")
        self.assert_np_arrays(model_item[3], np.array([4, 11]))

    def test_iter(self):
        for model_item in self.model:
            self.assertEqual(len(model_item), 4)
            self.assertEqual(model_item[0], "example.py")
            self.assertEqual(type(model_item[1]), Node)
            self.assertEqual(model_item[2],
                             "def f():\n    f2()\n    f2()\n    f3()\n    "
                             "f1()\n    f3()\n    f3()")
            self.assert_np_arrays(model_item[3], np.array([4, 11]))

    def test_positions_start(self):
        self.assert_np_arrays(self.model.positions_start, np.array([4]))

    def test_positions_end(self):
        self.assert_np_arrays(self.model.positions_end, np.array([11]))

    def test_positions(self):
        self.assert_np_arrays(self.model.positions, np.array([[4, 11]]))

    def test_bad_construct(self):
        repository = "repo_name"
        filenames = ["file"]
        sources = [""]
        uasts = [Node()]
        positions = np.array([[1, 2], [3, 4]])
        pos_start = [1, 3]
        pos_end = [2, 4]

        with self.assertRaises(ValueError):
            Snippet().construct(repository, filenames, sources, uasts)
        with self.assertRaises(ValueError):
            Snippet().construct(repository, filenames, sources, uasts,
                                pos_start, pos_end, positions)
        with self.assertRaises(ValueError):
            Snippet().construct(repository, filenames, sources, uasts,
                                positions_start=pos_start, positions=positions)
        with self.assertRaises(ValueError):
            Snippet().construct(repository, filenames, sources, uasts,
                                positions_end=pos_end, positions=positions)
        with self.assertRaises(ValueError):
            Snippet().construct(repository, filenames, sources, uasts,
                                positions_start=pos_start, positions_end=[0])

        Snippet().construct(repository, filenames, sources, uasts, pos_start, pos_end)
        Snippet().construct(repository, filenames, sources, uasts, positions=positions)


if __name__ == "__main__":
    unittest.main()
