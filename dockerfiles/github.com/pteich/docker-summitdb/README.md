# docker-summitdb
Dockerfile for SummitDB Docker image based on Alpine

[![Docker Build Statu](https://img.shields.io/docker/build/pteich/summitdb.svg)](https://hub.docker.com/r/pteich/summitdb/)

[SummitDB](https://github.com/tidwall/summitdb) Version 0.4.0

Small Docker image (24MB) for running a containerized [SummitDB](https://github.com/tidwall/summitdb).
SummitDB is an in-memory, NoSQL key/value database written in Go with a Redis-style API. It uses the Raft consensus to build a cluster.

SummitDB uses port 7481 per default for client and server-to-server communication.

This docker image makes use of gosu and dumb-init.

Built this container:

```
docker build --rm -t summitdb:latest ./
```

Run this container and use a local directory to persist data:

```
docker run -d -p 7481:7481 -v /path/to/local/dir:/data summitdb:latest
```

Run using [docker-compose](https://github.com/docker/compose):

```
docker-compose up -d
```
