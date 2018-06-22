import logging

from ast2vec.uast import UASTModel
from ast2vec.source import Source
from ast2vec.bblfsh_roles import CALL, CALL_CALLEE, FUNCTION_DECLARATION
from snippet_ranger.model2.base_split import Model2BaseSplit
from snippet_ranger.models.snippet import Snippet
from snippet_ranger.utils import uast_role_nodes, has_import, uast_to_bag, get_func_names_bow


class Source2Func(Model2BaseSplit):
    """
    This class splits files in source model to functions.
    It makes snippet from each function that satisfies the rule in `output_model_object_criteria`
    for every file that satisfies the rule in `input_model_object_criteria`.
    """
    MODEL_FROM_CLASS = Source
    MODEL_TO_CLASS = Snippet

    def __init__(self, libname, lib_funcs_bow=None, *args, **kwargs):
        """
        :param libname: Name of the library. All files without library usage are not handled.
        :param lib_funcs_bow: dictionary of function names (bag of words) which can be used from \
            the library. You can use just all function names or preprocess it somehow. For \
            example, remove functions that are common for language or internal functions.
        :param args: positional arguments to pass to :class:`Model2BaseSplit`.
        :param kwargs: key arguments to pass to :class:`Model2BaseSplit`.
        """
        super(Source2Func, self).__init__(*args, **kwargs)
        self.libname = libname
        self.lib_funcs_bow = lib_funcs_bow
        self.lib_funcs_bow_set = set(lib_funcs_bow)

        self._log.debug("lib_funcs_bow for {} lib is {}".format(libname, lib_funcs_bow))
        self.threshold = 0

    def input_model_object_criteria(self, model_object):
        """
        Filtration criteria for input model. Filter all files that are not use specified library.

        :param model_object: Object of Source model.
        :return: Library usage indicator.
        """
        filename, uast, source = model_object
        res = has_import(self.libname, uast)
        return res

    def split_model_object(self, model_from, model_object):
        """
        Parameters iterator for Function models initialisation.
        This models are split from

        :param model_from: Full Source model.
        :param model_object: Current Source model object.
        :return: parameters for :class:`Snippet` model __init__.
        """
        filename, uast, source = model_object
        have_funcs = False
        func_nodes = uast_role_nodes(uast, FUNCTION_DECLARATION)
        for func_node in func_nodes:
            have_funcs = True
            pos_start, pos_end = func_node.start_position.line-1, func_node.end_position.line
            func_source = "\n".join(source.splitlines()[pos_start:pos_end])
            yield filename, func_node, func_source, pos_start, pos_end
        if not have_funcs:
            yield filename, uast, source, 0, source.count("\n")

    def output_model_object_criteria(self, model_object):
        """
        Output model object filter criteria. Filter all :class:`Snippet` model entries that has
        no calls from the library bag of words are not use the library.

        :param model_object: Source model object.
        :return: Library functions usage indicator.
        """
        func_node = model_object[1]
        func_names = uast_to_bag(func_node, role=CALL_CALLEE)

        common = self.lib_funcs_bow_set & func_names.keys()
        if common == 0 and self._log.isEnabledFor(logging.DEBUG):
            self._log.DEBUG("There is no common functions in func_{}-{}_{}".format(
                model_object[4], model_object[5], model_object[0]))
            for key in func_names:
                if key in model_object[2]:
                    self._log.DEBUG("But I can find {} in source!".format(key))
                    self._log.error("Something strange. "
                                    "Here is the source code file:\n{}".format(model_object[2]))
                    break
        if len(common) > self.threshold:
            return True
        return False

    def construct(self, model_from, result):
        return Source2Func.MODEL_TO_CLASS().construct(model_from.repository, *zip(*result))


def process_lib_functions(functions_bow):
    """
    Remove bad function names from function bag of words. Specific for Python.
    Removes internal functions and test functions.
    """
    bad_prefixes = ["_", "test"]
    clear_functions_bow = {}
    for k, v in functions_bow.items():
        for bad_prefix in bad_prefixes:
            if k.startswith(bad_prefix):
                break
        else:
            clear_functions_bow[k] = v

    return clear_functions_bow


def source2func_entry(args):
    lib_model = UASTModel().load(args.library_uast)
    functions_bow = get_func_names_bow(lib_model)
    functions_bow = process_lib_functions(functions_bow)

    converter = Source2Func(args.library_name, functions_bow, log_level=args.log_level,
                            num_processes=args.processes,
                            overwrite_existing=args.overwrite_existing)
    converter.convert(args.input, args.output, pattern=args.filter)
