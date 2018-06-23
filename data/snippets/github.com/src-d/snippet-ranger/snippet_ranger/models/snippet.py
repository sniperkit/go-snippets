from ast2vec import Source
import numpy as np


class Snippet(Source):
    """
    This model can store code snippets. In general, code snippet is any part of source code file.
    For example, function declaration is a code snippet. So, this class is the same as source model
    but have start and end line positions of snippet location in file.
    You can use :class:`Source2Function` transformer to create function snippets from source model.
    """

    NAME = "snippet"

    def construct(self, repository, filenames, uasts, sources,
                  positions_start=None, positions_end=None, positions=None):
        super(Snippet, self).construct(repository=repository, filenames=filenames,
                                       sources=sources, uasts=uasts)
        if ((positions_start is None) ^ (positions_end is None)) or \
                ((positions_start is None) ^ (positions is not None)):
            raise ValueError("You should specify both positions_start and positions_end or "
                             "only a positions")
        if positions is None:
            if len(positions_start) != len(positions_end):
                raise ValueError("Length of positions_start ({}) and positions_end ({}) "
                                 "are not equal".format(len(positions_start), len(positions_end)))
            self._positions = np.array(list(zip(positions_start, positions_end)))
        else:
            self._positions = positions
        return self

    @property
    def names(self) -> list:
        """
        Creates the list of names for snippets in the model.
        Usually names are needed for topic modeling to represent different snippets as different
        documents. See `SnippetModel2BOW` transformer.
        """
        return ["{}/{}_{}_{}".format(self._repository, name, st, end).
                replace(":", "").replace(" ", "_")
                for name, (st, end) in zip(self._filenames, self._positions)]

    @property
    def positions(self):
        """
        Return start and end line positions of snippets.
        """
        return self._positions

    @property
    def positions_start(self):
        """
        Return start line position of snippets.
        """
        return self._positions[:, 0].T

    @property
    def positions_end(self):
        """
        Return end line position of snippets.
        """
        return self._positions[:, 1].T

    def __iter__(self):
        """
        Iterator over the items.
        """
        return zip(self._filenames, self._uasts, self._sources, self._positions)

    def __getitem__(self, item):
        """
        Returns file name, uast, source code and positions for the given snippet index.

        :param item: Snippet index.
        :return: file name, source code, uast, positions, where positions[0] is start and \
            positions[1] is end.
        """
        return super(Snippet, self).__getitem__(item) + (self._positions[item], )

    def _load_tree_kwargs(self, tree):
        tree_kwargs = super(Snippet, self)._load_tree_kwargs(tree)
        tree_kwargs["positions"] = np.array(tree["positions"])
        return tree_kwargs

    def _to_dict_to_save(self):
        save_dict = super(Snippet, self)._to_dict_to_save()
        save_dict["positions"] = self._positions
        return save_dict
