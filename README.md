# Debian

Base docker image based on Debian. Special variants for Sid / Jessie / Wheezy.

------

[![Docker Stars](https://img.shields.io/docker/stars/dockette/debian.svg?style=flat)](https://hub.docker.com/r/dockette/debian/)
[![Docker Pulls](https://img.shields.io/docker/pulls/dockette/debian.svg?style=flat)](https://hub.docker.com/r/dockette/debian/)

## Discussion / Help

[![Join the chat](https://img.shields.io/gitter/room/dockette/dockette.svg?style=flat-square)](https://gitter.im/dockette/dockette?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

## Image

- predefined user `dfx` with UID `1000`
- predefined constants:
    - `USER_UID`: `1000`
    - `USER_NAME`: `dfx`
    - `USER_HOME`: `/home/dfx`
- some optimalization for smaller image 

## Stretch

[![Docker Stars](https://img.shields.io/docker/stars/dockette/stretch.svg?style=flat)](https://hub.docker.com/r/dockette/stretch/)
[![Docker Pulls](https://img.shields.io/docker/pulls/dockette/stretch.svg?style=flat)](https://hub.docker.com/r/dockette/stretch/)

```
docker run -it --rm dockette/debian:stretch /bin/bash
```

```
docker run -it --rm dockette/stretch /bin/bash
```

## Jessie

[![Docker Stars](https://img.shields.io/docker/stars/dockette/jessie.svg?style=flat)](https://hub.docker.com/r/dockette/jessie/)
[![Docker Pulls](https://img.shields.io/docker/pulls/dockette/jessie.svg?style=flat)](https://hub.docker.com/r/dockette/jessie/)

```
docker run -it --rm dockette/debian:jessie /bin/bash
```

```
docker run -it --rm dockette/jessie /bin/bash
```

## Sid

[![Docker Stars](https://img.shields.io/docker/stars/dockette/sid.svg?style=flat)](https://hub.docker.com/r/dockette/sid/)
[![Docker Pulls](https://img.shields.io/docker/pulls/dockette/sid.svg?style=flat)](https://hub.docker.com/r/dockette/sid/)

```
docker run -it --rm dockette/debian:sid /bin/bash
```

```
docker run -it --rm dockette/sid /bin/bash
```

## Wheezy

[![Docker Stars](https://img.shields.io/docker/stars/dockette/wheezy.svg?style=flat)](https://hub.docker.com/r/dockette/wheezy/)
[![Docker Pulls](https://img.shields.io/docker/pulls/dockette/wheezy.svg?style=flat)](https://hub.docker.com/r/dockette/wheezy/)

```
docker run -it --rm dockette/debian:wheezy /bin/bash
```

```
docker run -it --rm dockette/wheezy /bin/bash
```
