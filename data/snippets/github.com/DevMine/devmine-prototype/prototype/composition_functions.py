"""This file provides different implementations for the composition function"""

# The name composition function is not very clear, actually. Should we rename
# it? Ranking function, score function, cool function...


def dummy(scores, query):
    """Always returns 0"""
    return 0


def sum_scores(scores, query):
    """Sums all the scores"""
    return sum(scores.values())  # self-documenting, indeed


def weighted_sum(scores, query):
    """Computes a weighted sum based on the query"""
    DEFAULT_WEIGHT = 0

    total = 0

    for feature in scores.keys():
        if feature in query:
            total += scores[feature] * query[feature]
        else:
            total += scores[feature] * DEFAULT_WEIGHT

    return total
