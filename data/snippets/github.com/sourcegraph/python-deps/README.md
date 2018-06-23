python-deps
===========

[![xrefs](https://sourcegraph.com/api/repos/github.com/sourcegraph/python-deps/badges/xrefs.png)](https://sourcegraph.com/github.com/sourcegraph/python-deps)
[![funcs](https://sourcegraph.com/api/repos/github.com/sourcegraph/python-deps/badges/funcs.png)](https://sourcegraph.com/github.com/sourcegraph/python-deps)
[![top func](https://sourcegraph.com/api/repos/github.com/sourcegraph/python-deps/badges/top-func.png)](https://sourcegraph.com/github.com/sourcegraph/python-deps)
[![library users](https://sourcegraph.com/api/repos/github.com/sourcegraph/python-deps/badges/library-users.png)](https://sourcegraph.com/github.com/sourcegraph/python-deps)
[![status](https://sourcegraph.com/api/repos/github.com/sourcegraph/python-deps/badges/status.png)](https://sourcegraph.com/github.com/sourcegraph/python-deps)

python-deps is a lightweight python module that computes the external
imports of a python project.  It does so statically, without importing
or running any code from the module itself.

The command and tests should be run from a clean virtualenv python.

This uses the standard lib associated with whatever python
installation used to run the module. A future version might support
specifying a path to a python env directory, to decouple the choice of
python used to execute the module vs. the python from which the names
of stdlib modules are detected.

Note: this module is a work-in-progress.

Installation
==========
```
pip install git+https://github.com/sourcegraph/python-deps
```

Running
=======
```
depdump.py path/to/project/root
```

Running tests
=============
- checkout: `git clone https://github.com/sourcegraph/python-deps.git`
- run: `nosetests -s`

Contributors
* Beyang Liu <beyang@sourcegraph.com>
