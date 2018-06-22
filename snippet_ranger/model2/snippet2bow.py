import os
from collections import defaultdict
import logging
import math

from ast2vec.bow import BOW
from ast2vec.model2.source2bow import UastModel2BOW, Uasts2BOW
from ast2vec.df import DocumentFrequencies
from ast2vec.token_parser import NoTokenParser
from ast2vec.uast_ids_to_bag import UastIds2Bag
from ast2vec.bblfsh_roles import SIMPLE_IDENTIFIER, CALL_CALLEE
import numpy
from scipy.sparse import csr_matrix

from snippet_ranger.models.snippet import Snippet


class Snippet2BOW(Uasts2BOW):
    def __init__(self, vocabulary: dict, docfreq: DocumentFrequencies,
                 getter: callable, token_parser=None, role=SIMPLE_IDENTIFIER):
        super().__init__(vocabulary, docfreq, getter)
        if token_parser is not None:
            self._uast2bag = UastIds2Bag(vocabulary, token_parser)
        self.role = role

    def __call__(self, file_uast):
        freqs = defaultdict(int)
        bag = self._uast2bag.uast_to_bag(self._getter(file_uast), self.role)
        for key, freq in bag.items():
            freqs[key] += freq
        missing = []
        for key, val in freqs.items():
            try:
                freqs[key] = math.log(1 + val) * math.log(
                    self._docfreq.docs / self._docfreq[self._reverse_vocabulary[key]])
            except KeyError:
                missing.append(key)
        for key in missing:
            del freqs[key]
        return dict(freqs)


class SnippetModel2BOW(UastModel2BOW):
    MODEL_FROM_CLASS = Snippet

    def __init__(self, topn, docfreq, *args, **kwargs):
        super().__init__(topn, docfreq, *args, **kwargs)
        vocabulary = {t: i for i, t in enumerate(self._tokens)}
        self._uasts2bow = Snippet2BOW(vocabulary, docfreq, lambda x: x)

    def convert_model(self, model: Snippet) -> BOW:
        bags = [self._uasts2bow(uast) for uast in model.uasts]
        data = list(zip(*[(bag_[x], i, x) for i, bag_ in enumerate(bags) for x in bag_]))
        matrix = csr_matrix((data[0], (data[1], data[2])),
                            shape=(len(bags), len(self._uasts2bow.vocabulary)),
                            dtype=numpy.float32)
        bow = BOW(log_level=logging.WARNING)
        bow.construct(repos=model.names, matrix=matrix, tokens=self._tokens)
        bow.meta["dependencies"] = [self._uasts2bow.docfreq]
        return bow


class SnippetModel2FuncCallsBOW(UastModel2BOW):
    # TODO(zurk): docstring

    MODEL_FROM_CLASS = Snippet

    def __init__(self, topn, docfreq, *args, **kwargs):
        super().__init__(topn, docfreq, *args, **kwargs)
        vocabulary = {t: i for i, t in enumerate(self._tokens)}
        self._uasts2bow = Snippet2BOW(vocabulary, docfreq, lambda x: x, NoTokenParser(),
                                      CALL_CALLEE)

    def convert_model(self, model: Snippet) -> BOW:
        bags = [self._uasts2bow(uast) for uast in model.uasts]
        data = list(zip(*[(bag_[x], i, x) for i, bag_ in enumerate(bags) for x in bag_]))
        matrix = csr_matrix((data[0], (data[1], data[2])),
                            shape=(len(bags), len(self._uasts2bow.vocabulary)),
                            dtype=numpy.float32)
        bow = BOW(log_level=logging.WARNING)
        bow.construct(repos=model.names, matrix=matrix, tokens=self._tokens)
        bow.meta["dependencies"] = [self._uasts2bow.docfreq]
        return bow


def snippet2bow_entry(args):
    df = DocumentFrequencies().load(args.docfreq)
    os.makedirs(args.output, exist_ok=True)
    converter = SnippetModel2BOW(args.vocabulary_size, df, num_processes=args.processes,
                                 overwrite_existing=args.overwrite_existing)
    converter.convert(args.input, args.output, pattern=args.filter)


def snippet2fc_bow_entry(args):
    df = DocumentFrequencies().load(args.docfreq)
    os.makedirs(args.output, exist_ok=True)
    converter = SnippetModel2FuncCallsBOW(args.vocabulary_size, df, num_processes=args.processes,
                                          overwrite_existing=args.overwrite_existing)
    converter.convert(args.input, args.output, pattern=args.filter)
