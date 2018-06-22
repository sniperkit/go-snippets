import ast
import sys
import os.path
from glob import glob
from collections import defaultdict

class Node(object):
    """Tree structure for managing python dependencies"""

    @classmethod
    def new_nonleaf(cls):
        return Node(False, defaultdict(Node.new_nonleaf))

    def __init__(self, isleaf, children_dict):
        self.children = children_dict
        self.isleaf = isleaf

    def __key(self):
        return (self.children, self.isleaf)

    def __eq__(self, other):
        return self.__key() == other.__key()

    def __hash__(self):
        return hash(self.__key())

    def to_dict(self):
        return {
            'children': { k: v.to_dict() for k, v in self.children.iteritems() },
            'isleaf': str(self.isleaf),
        }

    def leaves(self):
        for k, child in self.children.iteritems():
            if child.isleaf:
                yield k
            for l in child.leaves():
                yield k + '.' + l

    def add_path(self, path):   # absolute path.in.this.form
        if len(path) == 0: return
        components = path.split('.')
        node = self
        for component in components:
            node = node.children[component]
        node.isleaf = True  # set leaf

    def remove_path_and_children(self, path):
        if len(path) == 0: return
        components = path.split('.')
        prev = None
        node = self
        for component in components:
            prev = node
            node = node.children[component]
        if components[-1] in prev.children:
            del prev.children[components[-1]]

    def contains_prefix_of(self, path):
        components = path.split('.')
        prev = None
        node = self
        for component in components:
            if node.isleaf: return True
            if component in node.children:
                prev = node
                node = node.children[component]
            else:
                break
        return node.isleaf

    def print_tree(self):
        self.print_tree_prefix('')

    def print_tree_prefix(self, prefix):
        for k, v in self.children.iteritems():
            if v.isleaf is True:
                print '%s%s [leaf]' % (prefix, k)
            else:
                print '%s%s' % (prefix, k)
            v.print_tree_prefix('  ' + prefix)

class DepVisitor(ast.NodeVisitor):
    def __init__(self):
        self.imports = set([])

    def visit_Import(self, node):
        for alias in node.names:
            self.imports.add(alias.name)

    def visit_ImportFrom(self, node):
        if node.level > 0:
            return              # ignore relative imports

        for alias in node.names:
            if alias.name == '*':
                if node.module is not None:
                    self.imports.add(node.module)
            else:
                if node.module is not None:
                    impt = '%s.%s' % (node.module, alias.name)
                    self.imports.add(impt)

def __modules_with_root_module_path(path):
    """
    Returns all modules beneath the root module path. This treats all
    directories as packages regardless of whether or not they include
    a __init__.py.
    """
    modules = []
    if os.path.isfile(path) and os.path.splitext(path)[1] == '.py' and os.path.basename(path) != '__init__.py':
        name = os.path.splitext(os.path.basename(path))[0]
        modules.append(name)
    elif os.path.isdir(path):
        pkg_name = os.path.basename(path)
        modules.append(pkg_name)
        for ff in os.listdir(path):
            modules.extend(['.'.join([pkg_name, m]) for m in __modules_with_root_module_path(os.path.join(path, ff))])
    return modules

def paths_to_root_modules(rootpath, ignore_paths=[], followlinks=True):
    """
    Returns list of all paths to top-level (root) modules beneath
    rootpath. Optional arguments: follow symbolic links, list of
    directory paths to ignore (won't return any modules at or under
    this path).
    """

    if any([os.path.normpath(rootpath).startswith(os.path.normpath(ignore_path))
            for ignore_path in ignore_paths]):
        return []

    if os.path.isfile(rootpath) and os.path.splitext(rootpath)[1] == '.py':
        if rootpath.endswith('/setup.py') or rootpath == 'setup.py':
            return []
        else:
            return [rootpath]
    if os.path.exists(os.path.join(rootpath, '__init__.py')):
        return [rootpath]
    if os.path.isfile(rootpath) or (os.path.islink(rootpath) and not followlinks):
        return []

    module_paths = []
    for ff in os.listdir(rootpath):
        subpath = os.path.join(rootpath, ff)
        module_paths.extend(paths_to_root_modules(subpath, ignore_paths, followlinks))
    return module_paths

def modules_defined_in(path, ignore_paths=[], followlinks=True):
    rootpaths = paths_to_root_modules(path, ignore_paths, followlinks)
    modules = []
    for r in rootpaths:
        modules.extend(__modules_with_root_module_path(r))
    return modules

def root_modules_defined_in(path, ignore_paths=[], followlinks=True):
    """
    Paths passed as arguments should be absolute paths (there is no
    input checking).
    """
    rootpaths = paths_to_root_modules(path, ignore_paths, followlinks)
    rootmodules = []
    for r in rootpaths:
        rootmodules.append(os.path.splitext(os.path.basename(r))[0])
    return rootmodules

def import_tree_for_project(projectroot, **kwargs):
    """
    Provides tree of imports for the project. By default, ignores
    stdlib modules and internal modules. Also ignores explicit
    relative import paths (even when ignore_internal is False) and
    does not handle implicit relative import paths (it treats these as
    absolute paths).
    """

    ignore_stdlib = kwargs.get('ignore_stdlib', True)
    ignore_internal = kwargs.get('ignore_internal', True)

    ignore_tree = Node.new_nonleaf()
    if ignore_internal:
        for m in modules_defined_in(projectroot):
            ignore_tree.add_path(m)
    if ignore_stdlib:
        for m in stdlib_root_modules():
            ignore_tree.add_path(m)

    import_tree = Node.new_nonleaf()
    root_module_paths = paths_to_root_modules(projectroot)

    for root_module_path in root_module_paths:
        if os.path.isdir(root_module_path):
            pyfiles = py_files_in_dir(root_module_path)
            for pyfile in pyfiles:
                add_imports_for_file_to_tree(root_module_path, pyfile, import_tree, ignore_tree)
        else:
            add_imports_for_file_to_tree(root_module_path, root_module_path, import_tree, ignore_tree)
    return import_tree

def add_imports_for_file_to_tree(root_module_path, filename, import_tree, ignore_tree):
    """
    root_module_path is either a *.py file or a directory containing __init__.py
    """
    with open(filename) as ff:
        try:
            root = ast.parse(ff.read())
        except:
            sys.stderr.write('Could not parse file %s\n' % filename)
            return
        visitor = DepVisitor()
        visitor.visit(root)
        for impt in visitor.imports:
            if len(impt) == 0: continue # empty import
            if ignore_tree.contains_prefix_of(impt): continue # absolute path in ignore_tree
            # TODO(bliu): ignore implicit relative imports

            import_tree.add_path(impt)

def py_files_in_dir(rootdir, followlinks=True):
    for root, dirs, files in os.walk(rootdir, followlinks=followlinks):
        for ff in files:
            if os.path.splitext(ff)[1] == '.py':
                yield os.path.join(root, ff)

def stdlib_root_modules():
    """
    Finds stdlib python packages (packages that shouldn't be
    downloaded via pip.
    """
    stdlib_dir, sitepkg_dir, global_sitepkg_dir = python_stdlib_dirs()

    # python modules
    py_modules = root_modules_defined_in(stdlib_dir, [global_sitepkg_dir, sitepkg_dir])

    # c modules
    dynload_dir = os.path.join(stdlib_dir, 'lib-dynload/*')
    so_modules = [os.path.splitext(os.path.basename(path))[0] for path in glob(dynload_dir)]

    return set(so_modules) | set(py_modules) | set(sys.builtin_module_names)

def python_stdlib_dirs():
    """
    Returns (<stdlib-dir>, <sitepkg-dir>, <global-sitepkg-dir>). This
    exists because sysconfig.get_python_lib(standard_lib=True) returns
    something surprising when running a virtualenv python (the path to
    the global python standard lib directory, rather than the
    virtualenv standard lib directory), whereas
    sysconfig.get_python_lib(standard_lib=False) returns a path to the
    local site-packages directory. When processing the standard lib
    directory (global), we should ignore the global site-packages
    directory, not just the local one (which wouldn't get processed
    anyway).
    """
    import distutils.sysconfig as sysconfig
    sitepkg_dir = sysconfig.get_python_lib(standard_lib=False)
    stdlib_dir = sysconfig.get_python_lib(standard_lib=True)
    return (stdlib_dir, sitepkg_dir, os.path.join(stdlib_dir, 'site-packages'))
