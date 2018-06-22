import unittest
import os.path
from astdep import *
import setupdep
import json

testdir = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'testdata')
repodir_ast = os.path.join(testdir, "simplerepo")
bgendata = False                # True to regenerate test data

class TestASTDependencies(unittest.TestCase):
    def setUp(self):
        self.maxDiff = None

    def __check(self, expfilename, actual, overwrite):
        if overwrite:
            print 'Overwriting %s' % expfilename
            with open(expfilename, 'w') as expfile:
                expfile.write(actual)
            self.assertTrue(False)
        else:
            with open(expfilename) as expfile:
                self.assertMultiLineEqual(expfile.read(), actual)

    def test_import_tree_creation(self):
        expfilename = os.path.join(repodir_ast, 'import_tree.exp.json')
        deps = import_tree_for_project(repodir_ast, ignore_internal=False, ignore_stdlib=False)
        self.__check(expfilename, json.dumps(deps.to_dict(), indent=4, sort_keys=True), bgendata)

    def test_leaves(self):
        expfilename = os.path.join(repodir_ast, 'imports.exp.txt')
        deps = import_tree_for_project(repodir_ast, ignore_internal=False, ignore_stdlib=False)
        actual = '\n'.join(deps.leaves())
        self.__check(expfilename, actual, bgendata)

    def test_root_modules_defined_in(self):
        expfilename = os.path.join(repodir_ast, 'root_modules.exp.txt')
        actual = '\n'.join(root_modules_defined_in(repodir_ast))
        self.__check(expfilename, actual, bgendata)

    def test_modules_defined_in(self):
        expfilename = os.path.join(repodir_ast, 'modules.exp.txt')
        actual = '\n'.join(modules_defined_in(repodir_ast))
        self.__check(expfilename, actual, bgendata)

    def test_external_import_tree_stdlib(self):
        expfilename = os.path.join(repodir_ast, 'ext_import_tree_stdlib.exp.json')
        ext_deps = import_tree_for_project(repodir_ast, ignore_internal=True, ignore_stdlib=False)
        self.__check(expfilename, json.dumps(ext_deps.to_dict(), indent=4, sort_keys=True), bgendata)

    def test_external_import_tree_nostdlib(self):
        expfilename = os.path.join(repodir_ast, 'ext_import_tree_nostdlib.exp.json')
        ext_deps = import_tree_for_project(repodir_ast, ignore_internal=True, ignore_stdlib=True)
        self.__check(expfilename, json.dumps(ext_deps.to_dict(), indent=4, sort_keys=True), bgendata)

    def test_python_version(self):
        # Because this library looks through the standard lib, test
        # that the version of python is correct.  Only tested in
        # Python 2.7.5
        import sys
        self.assertTrue(sys.version_info > (2, 7))

class TestSetupPyDependencies(unittest.TestCase):
    def test_easy(self):
        repodir = os.path.join(testdir, "setup_repo_easy")
        self.assertTrue(setupdep.deps(repodir) == ['nose==1.3.0', 'wsgiref==0.1.2'])

    def test_hard(self):
        repodir = os.path.join(testdir, "setup_repo_hard")
        self.assertTrue(setupdep.deps(repodir) == ['nose==1.3.0', 'wsgiref==0.1.2'])

    def test_too_hard(self):
        repodir = os.path.join(testdir, "setup_repo_too_hard")
        self.assertTrue(setupdep.deps(repodir) is None)

    def test_setup_with_ast_method(self):
        repodir = os.path.join(testdir, "setup_repo_easy")
        deps = import_tree_for_project(repodir, ignore_stdlib=True, ignore_internal=True)
        self.assertTrue(len(deps.children) == 0)
