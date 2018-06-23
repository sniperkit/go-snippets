# RedPen Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/redpen.svg)
![](https://img.shields.io/docker/build/qmu1/redpen.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/redpen.svg)](https://microbadger.com/images/qmu1/redpen)

An alipine-based dockerfile which runs [RedPen](https://github.com/redpen-cc/redpen) (an open source proofreading tool) with python 2.7.

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/redpen/Dockerfile)

## Usage

```
docker run -it --rm -v $(pwd):/data qmu1/redpen:latest /bin/sh -c "redpen -c redpen-conf-ja.xml target_doc.txt -r json2"
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/redpen/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
