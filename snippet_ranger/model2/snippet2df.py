import tempfile
from typing import Union

from ast2vec.model2.source2df import Uast2DocFreq, MergeDocFreq
from ast2vec.df import DocumentFrequencies
from ast2vec.bblfsh_roles import SIMPLE_IDENTIFIER, CALL_CALLEE
from ast2vec.source import UASTModel
from ast2vec.uast_ids_to_bag import UastIds2Bag
from snippet_ranger.utils import get_func_names_bow
from snippet_ranger.model2.source2func import process_lib_functions

from modelforge import Model

from snippet_ranger.models.snippet import Snippet


class FilterVocabulary:
    def __init__(self, filter_set):
        self.filter_set = filter_set

    def __getitem__(self, item):
        if item in self.filter_set:
            return item
        raise KeyError


class Snippet2DocFreq(Uast2DocFreq):
    MODEL_FROM_CLASS = Snippet
    MODEL_TO_CLASS = DocumentFrequencies

    def __init__(self, role: int=SIMPLE_IDENTIFIER, filter_set: set =None, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self._uast2bag = UastIds2Bag(FilterVocabulary(filter_set) if filter_set else None)
        self.role = role

    def convert_model(self, model: Model) -> Union[Model, None]:
        for uast in model.uasts:
            contained = set()
            for key in self._uast2bag.uast_to_bag(uast, role=self.role):
                contained.add(key)
            for word in contained:
                self._df[word] += 1
            self._docs += 1


def snippet2df_entry(args):
    converter = Snippet2DocFreq(num_processes=args.processes,
                                overwrite_existing=args.overwrite_existing)
    with tempfile.TemporaryDirectory(dir=args.tmpdir, prefix="snippet2df") as tmpdir:
        converter.convert(args.input, tmpdir, pattern=args.filter)
        joiner = MergeDocFreq(num_processes=1)
        joiner.convert(tmpdir, args.output,
                       pattern="%s*.asdf" % DocumentFrequencies.NAME)


def snippet2fc_df_entry(args):
    lib_model = UASTModel().load(args.library_uast)
    functions_bow = get_func_names_bow(lib_model)
    functions_bow = process_lib_functions(functions_bow)

    converter = Snippet2DocFreq(num_processes=args.processes,
                                overwrite_existing=args.overwrite_existing,
                                role=CALL_CALLEE, filter_set=functions_bow)
    with tempfile.TemporaryDirectory(dir=args.tmpdir, prefix="snippet2df") as tmpdir:
        converter.convert(args.input, tmpdir, pattern=args.filter)
        joiner = MergeDocFreq(num_processes=1)
        joiner.convert(tmpdir, args.output,
                       pattern="%s*.asdf" % DocumentFrequencies.NAME)
