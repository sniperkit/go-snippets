# Base Image of Project Nutshells 🌰

[![Project Nutshells](https://img.shields.io/badge/Project-_Nutshells_🌰-orange.svg?maxAge=2592000)](https://github.com/quchao/nutshells/) [![Docker Build Build Status](https://img.shields.io/docker/build/nutshells/base.svg?maxAge=3600&label=Build%20Status)](https://hub.docker.com/r/nutshells/base/) [![Alpine Based](https://img.shields.io/badge/Alpine-3.7-0D597F.svg?maxAge=2592000)](https://alpinelinux.org/) [![CentOS Based](https://img.shields.io/badge/CentOS-7.5-932279.svg?maxAge=2592000)](https://www.centos.org/) [![MIT License](https://img.shields.io/github/license/quchao/nutshells.svg?maxAge=2592000&label=License)](https://github.com/quchao/nutshells/blob/master/LICENSE/)


## Variants:

| Tag | Description | 🐳 |
|:-- |:-- |:--:|
| `:latest` | Based on `alpine:latest`, includes a handy tool named [su-exec](https://github.com/ncopa/su-exec/) to exec cmds with different privileges in containers. | [![Dockerfile](https://img.shields.io/badge/Dockerfile-latest-22B8EB.svg?maxAge=2592000&style=flat-square)](https://github.com/QuChao/nutshells/blob/master/base/Dockerfile/) |
| `:edge` | Based on `alpine:edge`. | [![Dockerfile](https://img.shields.io/badge/Dockerfile-edge-22B8EB.svg?maxAge=2592000&style=flat-square)](https://github.com/QuChao/nutshells/blob/master/base/Dockerfile-edge/) |
| `:centos` | Based on `centos:latest` with [EPEL](https://fedoraproject.org/wiki/EPEL) & [gosu](https://github.com/tianon/gosu/) pre-installed. | [![Dockerfile](https://img.shields.io/badge/Dockerfile-centos-22B8EB.svg?maxAge=2592000&style=flat-square)](https://github.com/QuChao/nutshells/blob/master/base/Dockerfile-centos/) |


## Usage

### Customizing the image

#### By modifying the dockerfile

You may want to make some modifications to the image.
Pull the source code from GitHub, customize it, then build one by yourself:

> Get `<img_tag>` & `<dockerfile_name>` from [this sheet](#variants).

``` bash
git clone --depth 1 https://github.com/quchao/nutshells.git
docker image build \
       -q=false --rm=true --no-cache=true \
       -t nutshells/base<img_tag> \
       -f ./base/<dockerfile_name> \
       ./base
```

#### By committing the changes on a container

Otherwise just pull the image from the official registry, start a container and get a shell to it, [commit the changes](https://docs.docker.com/engine/reference/commandline/commit/) afterwards.

``` bash
docker container commit --change "Commit msg" <container_name> nutshells/base
```


## Contributing

[![Github Starts](https://img.shields.io/github/stars/quchao/nutshells.svg?maxAge=3600&style=social&label=Star)](https://github.com/quchao/nutshells/) [![Twitter Followers](https://img.shields.io/twitter/follow/chappell.svg?maxAge=3600&style=social&label=Follow)](https://twitter.com/chappell/)

> Follow GitHub's [*How-to*](https://opensource.guide/how-to-contribute/) guide for the basis.

Contributions are always welcome in many ways:

- Give a star to show your fondness;
- File an [issue](https://github.com/quchao/nutshells/issues) if you have a question or an idea;
- Fork this repo and submit a [PR](https://github.com/quchao/nutshells/pulls);
- Improve the documentation.


## Todo

- [x] Extract common code of entrypoint script into a shared lib.
- [x] Rewrite the examples according to the new sub-commands since docker 1.13+.

## Acknowledgments & Licenses

Unless specified, all codes of **Project Nutshells** are released under the [MIT License](https://github.com/quchao/nutshells/blob/master/LICENSE).

Other relevant softwares:

| Ware/Lib | License |
|:-- |:--:|
| [Docker](https://www.docker.com/) | [![License](https://img.shields.io/github/license/moby/moby.svg?maxAge=2592000&label=License)](https://github.com/moby/moby/blob/master/LICENSE/) |
| [Official Alpine Docker image](https://github.com/gliderlabs/docker-alpine) | [![License](https://img.shields.io/github/license/gliderlabs/docker-alpine.svg?maxAge=2592000&label=License)](https://github.com/gliderlabs/docker-alpine/blob/master/LICENSE/) |
| [Official CentOS Docker image](https://github.com/CentOS/sig-cloud-instance-images) | ![License](https://img.shields.io/badge/License-GNU_General_Public_License_v2.0-blue.svg?maxAge=2592000) |
| [su-exec](https://github.com/ncopa/su-exec/) | [![License](https://img.shields.io/github/license/ncopa/su-exec.svg?maxAge=2592000&label=License)](https://github.com/ncopa/su-exec/blob/master/LICENSE) |
| [EPEL](https://fedoraproject.org/wiki/EPEL) | [![License](https://img.shields.io/badge/License-Various-blue.svg?maxAge=2592000)](https://fedoraproject.org/wiki/Licensing:Main?rd=Licensing) |
| [gosu](https://github.com/tianon/gosu/) | [![License](https://img.shields.io/github/license/tianon/gosu.svg?maxAge=2592000&label=License)](https://github.com/tianon/gosu/blob/master/LICENSE/) |
