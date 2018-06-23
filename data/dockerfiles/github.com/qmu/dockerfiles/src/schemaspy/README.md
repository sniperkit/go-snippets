# SchemaSpy Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/schemaspy.svg)
![](https://img.shields.io/docker/build/qmu1/schemaspy.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/schemaspy.svg)](https://microbadger.com/images/qmu1/schemaspy)

An alipine-based dockerfile which runs [SchemaSpy v6 RC1](https://github.com/schemaspy/schemaspy) with MySQL Driver.

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/schemaspy/Dockerfile)

## Usage

```
docker run -it --rm \
    -v /path/to/export:/home/user/workdir \
    -v /path/to/config:/home/user/conf \
    --link mysql_container:mysql \
    qmu1/schemaspy \
    -configFile /home/user/conf/config.file -rails -hq
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/schemaspy/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
