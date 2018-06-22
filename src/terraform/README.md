# Terraform Dockerfile 

![](https://img.shields.io/docker/pulls/qmu1/terraform.svg)
![](https://img.shields.io/docker/build/qmu1/terraform.svg)
[![layers badge](https://images.microbadger.com/badges/image/qmu1/terraform.svg)](https://microbadger.com/images/qmu1/terraform)

An alipine-based dockerfile which runs [Terraform](https://www.terraform.io/).

## Supported tags and respective Dockerfile links

* `latest` [(Dockerfile)](https://github.com/qmu/dockerfiles/blob/master/src/terraform/Dockerfile)


## Usage

```
docker run --rm -it \
   qmu1/terraform /bin/sh -c \
   'terraform -v'
```

## DockerHub

Available on the DockerHub.

https://hub.docker.com/r/qmu1/terraform/

## License 

**MIT License**

Copyright (c) qmu Co.,Inc.
