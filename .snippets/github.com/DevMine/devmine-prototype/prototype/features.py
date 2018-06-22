"""This file provides abstraction over the set of features"""

from abc import ABCMeta, abstractmethod


class AbstractFeature(metaclass=ABCMeta):

    def __init__(self, data_source):
        self.data_source = data_source
        self.__verify_sufficient_data__(data_source)

    @abstractmethod
    def compute_score(self, entity):
        """Compute the score for this feature for the given entity"""
        pass

    @abstractmethod
    def __verify_sufficient_data__(self, data_source):
        """Make sure the data source supports everything this feature needs.
        If that is not the case, this method should raise a
        DataNotSufficientError"""
        pass

    @abstractmethod
    def __str__(self):
        pass


class DataNotSufficientError(Exception):

    def __init__(self, value):
        self.value = value

    def __str__(self):
        return repr(self.value)


class FollowersFeature(AbstractFeature):

    __key__ = 'followers'

    def compute_score(self, entity):
        return entity[self.__key__]

    def __verify_sufficient_data__(self, data_source):
        if (not data_source.developers
                or self.__key__ not in data_source.developers[0]):
            raise DataNotSufficientError("Either " + data_source +
                                         " is empty or '" + self.__key__ +
                                         "' was not valid as a key")
