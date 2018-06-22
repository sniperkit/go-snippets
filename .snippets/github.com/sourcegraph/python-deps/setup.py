from setuptools import setup

setup(
    name='astdep',
    version='0.1',
    url='http://github.com/sourcegraph/python-deps',
    py_modules=['astdep', 'setupdep'],
    scripts=['depdump.py', 'aliases.py'],
    author='Beyang Liu',
    description='a lightweight python module that statically '
                'computes the external imports of a python project',
    zip_safe=False,
    install_requires=[
        'nose==1.3.0',
        'wsgiref==0.1.2',
    ],
)
