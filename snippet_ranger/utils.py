from collections import defaultdict

from ast2vec import Source
from ast2vec import UASTModel
from ast2vec import bblfsh_roles


def uast_to_bag(uast, bag=None, role=bblfsh_roles.SIMPLE_IDENTIFIER):
    """
    Convert UAST to bag of words for certain role.

    :param uast: Uast to process.
    :param bag: Specify existing bag of words if you want to update it.
    :param role: Role or list of roles to get from uast and convert to bag of words.
    :return: bag of words for tokens of nodes.
    """
    if bag is None:
        bag = defaultdict(int)
    # TODO(zurk): Can be optimised if we need. uast_role_nodes make a list from one role.
    for node in uast_role_nodes(uast, roles=role):
        if node.token != "":
            bag[node.token] += 1
    return bag


def uast_role_nodes(uast, roles=None) -> iter:
    """
    Filter UAST by provided roles and iterate trough corresponding nodes of uast.
    It is a generator.

    :param uast:
    :param roles: roles to filter.
    :return: iterator trough corresponding nodes.
    """
    if roles is None:
        roles = [bblfsh_roles.SIMPLE_IDENTIFIER]
    elif isinstance(roles, int):
        roles = [roles]
    elif not isinstance(roles, list):
        raise TypeError()

    stack = [uast]
    while stack:
        node = stack.pop(0)
        for role_id in roles:
            if role_id in node.roles:
                yield node
                break
        stack.extend(node.children)


def _iter_imports(uast):
    """
    Internal helper function to iterate through uast's imports
    """
    for n in uast_role_nodes(uast, [bblfsh_roles.IMPORT_PATH, bblfsh_roles.IMPORT_ALIAS]):
        for x in n.token.split("."):
            yield x


def get_imports(uast):
    """
    Return all module imports from the uast.
    It is specific for Python because imports should be splited by dot.

    :param uast:
    :return: set of import names
    """
    return set(_iter_imports(uast))


def has_import(libname, uast) -> bool:
    """
    Check `libname` import in the uast

    :param libname: name of library to check
    :param uast:
    :return:
    """
    for imp in _iter_imports(uast):
        if imp == libname:
            return True
    return False


def get_func_names_bow(uast_model: UASTModel) -> dict:
    """
    Returns bag of words for functions in all Source model.

    :param uast_model: UAST model
    :return: dictionary with words as keys
    """
    func_class_names = None
    for record in uast_model:
        func_class_names = uast_to_bag(record[1], func_class_names,
                                       bblfsh_roles.FUNCTION_DECLARATION_NAME)
    return func_class_names
