# Composer Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/composer.svg)
![](https://img.shields.io/docker/build/qmu1/composer.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/composer.svg)](https://microbadger.com/images/qmu1/composer)

An alipine-based dockerfile which runs [Composer](https://github.com/composer/composer) with [prestissimo](https://github.com/hirak/prestissimo).

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/composer/Dockerfile)

## Usage

```
docker run -v $(pwd):/app --rm -ti qmu1/composer:latest /bin/sh -c "\
    composer -v \
"
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/composer/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
