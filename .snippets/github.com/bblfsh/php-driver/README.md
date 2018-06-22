# php-driver  ![Driver Status](https://img.shields.io/badge/status-beta-dbd25c.svg) [![Build Status](https://travis-ci.org/bblfsh/php-driver.svg?branch=master)](https://travis-ci.org/bblfsh/php-driver) ![Native Version](https://img.shields.io/badge/php%20version-7.1.17--r0-aa93ea.svg) ![Go Version](https://img.shields.io/badge/go%20version-1.9-63afbf.svg)

php driver for [babelfish](https://github.com/bblfsh/bblfshd).


Development Environment
-----------------------

Requirements:
- `docker`
- [`bblfsh-sdk`](https://github.com/bblfsh/sdk) _(go get -u gopkg.in/bblfsh/sdk.v2/...)_
- UAST converter dependencies _(go get -t -v ./...)_

To initialize the build system execute: `bblfsh-sdk prepare-build`, at the root of the project. This will install the SDK at `.sdk` for this driver.

To execute the tests just execute `make test`, this will execute the test over the native and the go components of the driver. Use `make test-native` to run the test only over the native component or `make test-driver` to run the test just over the go component.

The build is done executing `make build`. To evaluate the result using a docker container, execute:
`docker run -it bblfsh/php-driver:dev-<commit[:7]>-dirty`


License
-------

GPLv3, see [LICENSE](LICENSE)



