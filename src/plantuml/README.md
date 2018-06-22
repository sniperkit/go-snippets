# PlantUML Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/plantuml.svg)
![](https://img.shields.io/docker/build/qmu1/plantuml.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/plantuml.svg)](https://microbadger.com/images/qmu1/plantuml)

An alipine-based dockerfile which runs [PlantUML](https://github.com/plantuml/plantuml) with python.

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/plantuml/Dockerfile)

## Usage

```
docker run -it --rm \
    -v $(pwd):/workspace \
    qmu1/plantuml:latest \
    /workspace/path/to/uml/files
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/plantuml/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
