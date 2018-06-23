#!/usr/bin/env python2.7

import astdep
import setupdep

import argparse
import sys

def debug():
    import pdb, traceback, code, sys
    def handle_exception(tp, val, tb):  # open up debugger on exception
        traceback.print_tb(tb)
        sys.stderr.write('ERROR: %s' % str(val))
        pdb.post_mortem(tb)
    sys.excepthook = handle_exception

def get_args():
    parser = argparse.ArgumentParser(description='Grab dependencies from a setup.py file')
    parser.add_argument('rootdir', help='path to root project directory')
    parser.add_argument('--debug', action='store_true', help='enable if you want to enable debugging stacktrace')
    parser.add_argument('--method', type=str, default='any', help='Method of acquiring dependencies.  Either "ast", "setup.py", or "any"')
    parser.add_argument('--toplevel', action='store_true', help='Used in conjunction with --method "ast".  Enable if you want to print just top-level modules')
    return parser.parse_args()

def print_ast_deps(toplevel):
    deps = astdep.import_tree_for_project(args.rootdir, ignore_stdlib=True, ignore_internal=True)
    if toplevel:
        for k in sorted(deps.children.keys()):
            print k
    else:
        deps.print_tree()

def print_setup_deps():
    deps = setupdep.deps(args.rootdir)
    if deps is None:
        return False
    for dep in deps:
        print dep
    return True

if __name__ == '__main__':
    args = get_args()
    if args.debug:
        sys.stderr.write('!!!!! DEBUG ENABLED !!!!!\n')
        debug()

    if args.method == 'ast':
        print_ast_deps(args.toplevel)
    elif args.method == 'setup.py':
        success = print_setup_deps()
        if not success:
            sys.stderr.write('Error: no setup.py file found\n')
            sys.exit(1)
    elif args.method == 'any':
        if print_setup_deps() is False:
            print_ast_deps(True)
    else:
        sys.stderr.write('Error due to unrecognized method: %s\n' % args.method)
        sys.exit(1)
