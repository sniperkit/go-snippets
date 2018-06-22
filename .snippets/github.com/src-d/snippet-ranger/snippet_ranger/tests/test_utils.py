import unittest

from ast2vec import Source
from ast2vec import bblfsh_roles

from snippet_ranger import utils
from snippet_ranger.tests import models


class UtilsTests(unittest.TestCase):
    def test_get_func_names_bow(self):
        source = Source().load(models.TEST_LIB)
        bow = utils.get_func_names_bow(source)
        true_bow = {
            "f1": 1,
            "f2": 1,
            "f3": 1,
            "f35": 1}
        self.assertEqual(bow, true_bow)

    def test_uast_to_bag(self):
        source = Source().load(models.TEST_REPO)
        uast = source.uasts[0]
        bag = utils.uast_to_bag(uast)
        true_bag = {
            "test_lib": 2,
            "f": 2,
            "f1": 2,
            "f2": 2,
            "f3": 4}
        self.assertEqual(bag, true_bag)

    def test_get_imports(self):
        source = Source().load(models.TEST_REPO)
        uast = source.uasts[0]
        imports = utils.get_imports(uast)
        true_imports = {"f1", "test_lib"}
        self.assertEqual(imports, true_imports)

        source = Source().load(models.TEST_LIB)
        uast = source.uasts[0]
        imports = utils.get_imports(uast)
        true_imports = set()
        self.assertEqual(imports, true_imports)

    def test_has_import(self):
        source = Source().load(models.TEST_REPO)
        self.assertTrue(utils.has_import("f1", source.uasts[0]))
        self.assertTrue(utils.has_import("test_lib", source.uasts[0]))

        source = Source().load(models.TEST_LIB)
        self.assertFalse(utils.has_import("f1", source.uasts[0]))
        self.assertFalse(utils.has_import("test_lib", source.uasts[0]))

if __name__ == "__main__":
    unittest.main()
