# Apache-PHP7 Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/apache-php7.svg)
![](https://img.shields.io/docker/build/qmu1/apache-php7.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/apache-php7.svg)](https://microbadger.com/images/qmu1/apache-php7)

An alipine-based dockerfile which runs [Apache with PHP7](http://www.apache-php7.org/).

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/apache-php7/Dockerfile)

## Usage

```
docker run --rm -v $(pwd):/project qmu1/apache-php7:latest /bin/sh -c "\
    apache-php7 build --clean -f /project/path/to/your/apache-php7.yml -d ${docs_dir}"
"
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/apache-php7/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
