# MailCatcher Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/mailcatcher.svg)
![](https://img.shields.io/docker/build/qmu1/mailcatcher.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/mailcatcher.svg)](https://microbadger.com/images/qmu1/mailcatcher)

An alipine-based dockerfile which runs [MailCatcher](https://github.com/sj26/mailcatcher).

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/mailcatcher/Dockerfile)

## Usage

```
docker run -d -p 1080:80 -p 1025:25 --name mail qmu1/mailcatcher
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/mailcatcher/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
