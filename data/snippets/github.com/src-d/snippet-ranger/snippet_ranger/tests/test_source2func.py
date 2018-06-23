import argparse
import logging
import os
import unittest
import tempfile

import asdf
from ast2vec import Source
from ast2vec.bblfsh_roles import Node
from modelforge import split_strings

from snippet_ranger.model2.source2func import Source2Func
from snippet_ranger.model2.source2func import process_lib_functions, source2func_entry
from snippet_ranger.models.snippet import Snippet
from snippet_ranger.tests import models
from snippet_ranger.utils import get_func_names_bow
from snippet_ranger.tests.test_snippet import validate_asdf_file


class Source2FuncTests(unittest.TestCase):
    def test_get_func_names_bow(self):
        in_bow = {
            "_xxx": 1,
            "xxx": 1,
            "_yyy": 1,
            "yyy": 1,
            "test_x": 1,
            "test_y": 1,
        }
        out_bow = process_lib_functions(in_bow)
        self.assertEqual(out_bow, {"xxx": 1, "yyy": 1})

    def test_convert_model(self):
        from ast2vec.bblfsh_roles import Node

        lib_model = Source().load(models.TEST_LIB)
        functions_bow = get_func_names_bow(lib_model)
        functions_bow = process_lib_functions(functions_bow)
        self.assertEqual(functions_bow, {"f1": 1, "f2": 1, "f3": 1, "f35": 1})

        repo_model = Source().load(models.TEST_REPO)

        converter = Source2Func(models.LIB_NAME, functions_bow)
        functons = converter.convert_model(repo_model)

        self.assertEqual(len(functons), 1)
        functon_obj = functons[0]
        self.assertEqual(functon_obj[0], "example.py")
        self.assertEqual(type(functon_obj[1]), Node)
        self.assertEqual(functon_obj[2],
                         "def f():\n    f2()\n    f2()\n    f3()\n    f1()\n    f3()\n    f3()")
        self.assertEqual(functon_obj[3][0], 4)
        self.assertEqual(functon_obj[3][1], 11)

    def test_source2func_object(self):
        with tempfile.TemporaryDirectory() as tmpdir:
            args = argparse.Namespace(library_uast=models.TEST_LIB,
                                      library_name=models.LIB_NAME,
                                      input=models.DATA_DIR,
                                      output=tmpdir,
                                      filter="**/*.asdf",
                                      log_level=logging.INFO,
                                      processes=1,
                                      overwrite_existing=True)
            source2func_entry(args)
            validate_asdf_file(self, os.path.join(tmpdir, "source_test_repo.asdf"))


if __name__ == "__main__":
    unittest.main()
