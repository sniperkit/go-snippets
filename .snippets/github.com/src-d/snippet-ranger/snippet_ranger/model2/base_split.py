from typing import Iterable

from modelforge import Model
from ast2vec import Model2Base


class Model2BaseSplit(Model2Base):
    """
    Class to split the model content and save it in a new model.
    `Model2Base.MODEL_FROM` should be iterable and iterator returns will be called model objects.
    For example, one can split objects of :class:`ast2vec.Source` to functions parts and store it
    as :class:`snippet_ranger.Snippet` model (see :class:`snippet_ranger.Source2Func` class).

    This class the next pipeline:
    1. Filter unrelated objects in the input Model with `input_model_object_criteria` function.
    2. Split one model object into several objects of new model with `split_model_object` function.
    3. Discard some new objects appeared on step 2 with `output_model_object_criteria` function.
    4. Save new model with all remaining object.

    Produce exactly one model from one.
    """
    def __init__(self, *args, **kwargs):
        super(Model2BaseSplit, self).__init__(*args, **kwargs)

    def filter_input_model_content(self, model) -> Iterable:
        """
        Filter object of input model by criteria in input_model_object_criteria
        :return: Remaining objects.
        """
        return filter(self.input_model_object_criteria, model)

    def input_model_object_criteria(self, model_object):
        """
        Criteria for input model filtering.
        No filtration by default.
        :param model_object: Model object.
        :return:
        """
        return True

    def filter_output_model_object(self, model_object) -> Iterable:
        """
        Filter object of output model by criteria in output_model_object_criteria
        :return:
        """
        return filter(self.output_model_object_criteria, model_object)

    def output_model_object_criteria(self, model_object):
        """
        Criteria for output objects filtering.
        No filtration by default.
        :param model_object: Model object.
        :return: Remaining objects.
        """
        return True

    def convert_model(self, model_from: Model) -> Model:
        result = []
        for model_object in self.filter_input_model_content(iter(model_from)):
            out_model_part = self.split_model_object(model_from, model_object)
            for out_model_object in self.filter_output_model_object(out_model_part):
                result.append(out_model_object)
        if result:
            out_model = self.construct(model_from, result)
        else:
            raise ValueError("All functions are filtered and you get empty model.")
        return out_model

    def split_model_object(self, model_from: Model, model_object) -> Iterable:
        """
        This function must be implemented in the child classes.

        :param model_from: The full input model instance.
        :param model_object: The current model object instance to split.
        :return: List or tuple of parameters for output model objects initialization.
        """
        raise NotImplementedError

    def construct(self, model_from: Model, result: list) -> Model:
        """
        Construct out model from objects.
        This function must be implemented in the child classes.

        :param model_from: The full input model instance.
        :param result: objects from `split_model_object`
        :return: Output model.
        """
        raise NotImplementedError
