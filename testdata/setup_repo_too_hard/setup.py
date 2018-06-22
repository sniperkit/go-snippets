from setuptools import setup

foo = ['one']
bar = ['two']

setup(
    name='foo',
    version='0.1',
    url='',
    py_modules=[],
    scripts=[],
    author='John Doe',
    description='yadda yadda',
    zip_safe=False,
    install_requires=foo+bar,
)
