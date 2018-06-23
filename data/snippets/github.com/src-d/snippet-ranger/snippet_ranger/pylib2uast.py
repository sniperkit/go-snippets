import os
import logging
from ast2vec.repo2.uast import repos2uast_entry
from ast2vec.enry import install_enry


def pylib2uast_entry(args):
    log = logging.getLogger("pylib2uast")
    module_names = args.input
    module_dirs = []
    for lib_name in args.input:
        try:
            module = __import__(lib_name)
            module_dir = os.path.abspath(os.path.dirname(module.__file__))
            module_dirs.append(module_dir)
        except ModuleNotFoundError as e:
            log.error("No module named '%s'. Skipping.")

    args.input = module_dirs

    install_enry()
    repos2uast_entry(args)
    for file in os.listdir(args.output):
        for name in module_names:
            if name in file:
                break
        os.rename(os.path.join(args.output, file), os.path.join(args.output, name + '.asdf'))
