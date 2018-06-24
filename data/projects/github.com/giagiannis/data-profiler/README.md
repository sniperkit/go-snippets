data-profiler [![Build Status](https://travis-ci.org/giagiannis/data-profiler.svg?branch=master)](https://travis-ci.org/giagiannis/data-profiler) [![goreport](https://goreportcard.com/badge/github.com/giagiannis/data-profiler)](https://goreportcard.com/report/github.com/giagiannis/data-profiler) [![Coverage Status](https://coveralls.io/repos/github/giagiannis/data-profiler/badge.svg?branch=master)](https://coveralls.io/github/giagiannis/data-profiler?branch=master) [![Docker Automated build](https://img.shields.io/docker/automated/jrottenberg/ffmpeg.svg)](https://hub.docker.com/r/ggian/data-profiler/)
=============
__data-profiler__ is a Go project used to transform a set of datasets, based on a set of characteristics (distribution similarity, correlation, etc.), in order to model the behavior of an operator, applied on top of them using Machine Learning techniques.


Screenshots
-----------

![Similarity Matrix](https://github.com/giagiannis/data-profiler/raw/master/_imgs/SM.png "Dataset Similarity Matrix")

![Dataset Space](https://github.com/giagiannis/data-profiler/raw/master/_imgs/DS.png "Dataset Space")

![SVM Modeling](https://github.com/giagiannis/data-profiler/raw/master/_imgs/SVMModeling.png "Operator Modeling with SVM")

![SVM Residuals Distribution](https://github.com/giagiannis/data-profiler/raw/master/_imgs/ResidualDistribution.png "SVM Residual Distribution")

Installation
------------
You have two ways of installing __data-profiler__:

1. Through Go:

```bash
# GOPATH must be set
~> go get github.com/giagiannis/data-profiler
```

2. Using Docker:

```bash
~> docker pull ggian/data-profiler
```

Usage
-----
__data-profiler__ can be used both through a CLI and a Web interface.

1. CLI 

You can access the CLI client through the __data-profiler-utils__ binary.

```bash
~> $GOPATH/bin/data-profiler-utils
```

This previous command will give an overview of the available actions.

__Note:__ use this client only if you know how data-profiler works.

2. Web UI

First run the Docker container, providing a directory with the dataset files. 

```bash
~> docker run -v /src/datasets:/datasets -p 8080:8080 -d ggian/data-profiler
```

This command mounts the host's _/src/datasets_ directory to the container and forwards the host's 8080 port to the container. After the successful start of the container, go to _http://dockerhost:8080_ and insert the first set of datasets for analysis.


License
-------
Apache License v2.0 (see [LICENSE](LICENSE) file for more)


Contact
-------
Giannis Giannakopoulos ggian@cslab.ece.ntua.gr
