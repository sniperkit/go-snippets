import test_lib
from test_lib import f1


def f():
    f2()
    f2()
    f3()
    f1()
    f3()
    f3()


def f3():
    f()
