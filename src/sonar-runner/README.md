# MkDocs Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/mkdocs.svg)
![](https://img.shields.io/docker/build/qmu1/mkdocs.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/mkdocs.svg)](https://microbadger.com/images/qmu1/mkdocs)

An alipine-based dockerfile which runs [MkDocs](http://www.mkdocs.org/) with its [material theme](http://squidfunk.github.io/mkdocs-material/).

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/sonar-runner/Dockerfile)

## Usage

```
docker run --rm -v $(pwd):/project qmu1/mkdocs:latest /bin/sh -c "\
    mkdocs build --clean -f /project/path/to/your/mkdocs.yml -d ${docs_dir}"
"
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/mkdocs/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
