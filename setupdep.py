"""
Retrieves dependencies from setup.py file
"""

import ast
from os import path

import pdb

class SetupVisitor(ast.NodeVisitor):
    def __init__(self):
        self.install_requires_parse_fail = False
        self.candidates = {}
        self.deps = []

    def visit_Assign(self, node):
        if len(node.targets) == 1 and isinstance(node.targets[0], ast.Name):
            name = node.targets[0].id
            if 'req' in name.lower() and isinstance(node.value, ast.List):
                self.candidates[name] = self.deps_from_ast_list(node.value)

    def visit_Call(self, node):
        if isinstance(node.func, ast.Name) and node.func.id == 'setup':
            for kw in node.keywords:
                if kw.arg == 'install_requires':
                    self.record_deps_from_expr(kw.value)

    def deps_from_ast_list(self, ast_list):
        d = []
        for elt in ast_list.elts:
            if isinstance(elt, ast.Str):
                d.append(elt.s)
        return d

    def record_deps_from_expr(self, expr_node):
        if isinstance(expr_node, ast.List):
            self.deps.extend(self.deps_from_ast_list(expr_node))
        elif isinstance(expr_node, ast.Name):
            if expr_node.id in self.candidates:
                self.deps.extend(self.candidates[expr_node.id])
        else:
            self.install_requires_parse_fail = True

def deps(project_dir):
    """Returns (deps, success)"""
    setup_file = path.join(project_dir, 'setup.py')
    if path.exists(setup_file):
        return deps_from_setup_file(setup_file)
    return None

def deps_from_setup_file(setup_file):
    with open(setup_file) as setupf:
        root = ast.parse(setupf.read())
        visitor = SetupVisitor()
        visitor.visit(root)
        if visitor.install_requires_parse_fail:
            return None
        else:
            return visitor.deps
    return None
