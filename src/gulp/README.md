# Gulp Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/gulp.svg)
![](https://img.shields.io/docker/build/qmu1/gulp.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/gulp.svg)](https://microbadger.com/images/qmu1/gulp)

An alipine-based dockerfile which runs [Gulp](https://github.com/gulpjs/gulp) and [Bower](https://github.com/bower/bower).

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/gulp/Dockerfile)

## Usage

```
docker run -it --rm -v $(pwd):/workspace qmu1/gulp:latest /bin/sh -c " \
    gulp your_task --cwd path/to/gulpfile/ \
"
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/gulp/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
