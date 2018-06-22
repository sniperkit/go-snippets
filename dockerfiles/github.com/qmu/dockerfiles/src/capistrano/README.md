# Capistrano Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/capistrano.svg)
![](https://img.shields.io/docker/build/qmu1/capistrano.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/capistrano.svg)](https://microbadger.com/images/qmu1/capistrano)

An alipine-based dockerfile which runs [Capistrano](https://github.com/capistrano/capistrano).

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/capistrano/Dockerfile)

## Usage

```
docker run -it --rm \
    -e "APP_SERVER_KEY=hoge.pem" \
    -e "GITHUB_KEY=fuga.pem" \
    -v ~/.ssh:/root/.ssh \
    -v $PROJECT_PATH:/root/source \
    qmu1/capistrano ./cap production deploy
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/capistrano/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
