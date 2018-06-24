#!/usr/bin/env python3

import ast
import argparse
import os
from os import path
import json

class FromImportAsVisitor(ast.NodeVisitor):
    def __init__(self, parent_module):
        self.imports = []
        self.parent_module = parent_module

    def visit_ImportFrom(self, node):
        if node.module is None:
            return

        if node.level == 0:
            module = node.module
        elif node.level == 1:
            module = self.parent_module + '.' + node.module
        else:
            # not handling relative imports with ..
            return

        aliases = {'from': module, 'imports': []}
        for alias in node.names:
            as_ = alias.asname
            if as_ is None: as_ = ''
            aliases['imports'].append({'name': alias.name, 'as': as_})
        self.imports.append(aliases)

def get_aliases(rootdir):
    rootdir = path.abspath(rootdir)
    top_module_dir = rootdir
    if path.exists(path.join(rootdir, '__init__.py')):
        top_module_dir = path.dirname(rootdir)

    module_aliases = []
    for dirpath, dirnames, filenames in os.walk(rootdir):
        for filename in filenames:
            if filename == '__init__.py':
                module_name = path.relpath(dirpath, top_module_dir).replace('/', '.')
                if module_name.startswith('..'):
                    module_name = module_name[2:]
                with open(path.join(dirpath, filename)) as f:
                    try: # best effort
                        root = ast.parse(f.read())
                        visitor = FromImportAsVisitor(module_name)
                        visitor.visit(root)
                        module_aliases.append({'module': module_name, 'reExports': visitor.imports})
                    except:
                        pass
    return module_aliases

if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Parses all "re-exports" or module aliases in a project and prints them in JSON')
    parser.add_argument('rootdir', help='path to root project directory')

    args = parser.parse_args()
    aliases_json = json.dumps(get_aliases(args.rootdir), separators=(',', ':'), indent=None, sort_keys=True)
    print aliases_json
