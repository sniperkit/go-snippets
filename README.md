# erlang-driver  ![Driver Status](https://img.shields.io/badge/status-planning-e08dd1.svg) [![Build Status](https://travis-ci.org/bblfsh/erlang-driver.svg?branch=master)](https://travis-ci.org/bblfsh/erlang-driver) ![Native Version](https://img.shields.io/badge/erlang%20version-erts--8.3-aa93ea.svg) ![Go Version](https://img.shields.io/badge/go%20version-1.8-63afbf.svg)

erlang driver for [babelfish](https://github.com/bblfsh/server).


Development Environment
-----------------------

Requirements:
- `docker`
- [`bblfsh-sdk`](https://github.com/bblfsh/sdk) _(go get -u github.com/bblfsh/sdk/...)_
- UAST converter dependencies _(go get -t -v ./...)_

To initialize the build system execute: `bblfsh-sdk prepare-build`, at the root of the project. This will install the SDK at `.sdk` for this driver.

To execute the tests just execute `make test`, this will execute the test over the native and the go components of the driver. Use `make test-native` to run the test only over the native component or `make test-driver` to run the test just over the go component.

The build is done executing `make build`. To evaluate the result using a docker container, execute:
`docker run -it bblfsh/erlang-driver:dev-<commit[:6]>`


License
-------

GPLv3, see [LICENSE](LICENSE)



