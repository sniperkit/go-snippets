# Severspec Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/serverspec.svg)
![](https://img.shields.io/docker/build/qmu1/serverspec.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/serverspec.svg)](https://microbadger.com/images/qmu1/serverspec)

An alipine-based dockerfile which runs [Serverspec](http://serverspec.org).

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/serverspec/Dockerfile)


## Usage

```
docker run --rm -it -u docker\
   -v ~/.ssh:/home/docker/.ssh \
   -v $(pwd):/home/docker/serverspec \
   qmu1/serverspec /bin/sh -c \
   'cd ~/serverspec;
   rake spec;
   '
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/serverspec/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
