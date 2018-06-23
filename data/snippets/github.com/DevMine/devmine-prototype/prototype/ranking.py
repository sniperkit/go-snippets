"""This file provides abstraction over the taks of computing the ranking"""


class Ranking:
    """Compute the ranking for the developers based on a given composition
    function"""

    def __init__(self, data_source, composition_function):
        self.data_source = data_source
        self.composition_function = composition_function

        self.scores = dict(data_source.scores)          # clone!
        self.features = list(data_source.features)      # clone!
        self.developers = list(data_source.developers)  # clone!

    def rank_all(self, query):
        g = self.composition_function
        scores = self.scores

        result = [(g(scores[dev], query), dev) for dev in self.developers]

        result.sort()
        result.reverse()

        return [(dev, score) for (score, dev) in result]
