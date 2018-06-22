# PlantUML (with Gulp) Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/plantuml-gulp.svg)
![](https://img.shields.io/docker/build/qmu1/plantuml-gulp.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/plantuml-gulp.svg)](https://microbadger.com/images/qmu1/plantuml-gulp)

An alipine-based dockerfile which runs [PlantUML](https://github.com/plantuml/plantuml-gulp), [Gulp](https://github.com/gulpjs/gulp) within .

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/plantuml-gulp/Dockerfile)

## Usage

```
docker run \
    -it \
    --rm \
    -v $(pwd)/uml/:/workspace/uml \
    -v $(pwd)/img/:/workspace/img \
    qmu1/plantuml-gulp:latest \
    /bin/sh -c "
        gulp watch
    "
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/plantuml-gulp/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
