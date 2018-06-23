# Containerisation with Docker

In this example we put every service in its own Docker container. We reuse the source code from [previous example](../03-consul). 

For each custom service (`Sensor`, `Worker`, `App`) we have to create Docker image first by writing [Dockerfile](./images/sensor/Dockerfile):

```
FROM gliderlabs/alpine:3.4
COPY sensor /bin
WORKDIR bin
ENTRYPOINT ["sensor"]
```

For all other containers (*nsqd*, *nsqlookupd*, *nsqadmin*, *consul*, *registrator*) we use images available on Docker hub.

### Requirements

- [ruby](https://github.com/rbenv/rbenv)
- [thor](https://github.com/erikhuda/thor)
- [golang](https://golang.org/doc/install)
- [go-nsq](https://github.com/nsqio/go-nsq)
- [docker](https://www.docker.com/)
- [docker-compose](https://docs.docker.com/compose/)

In this example we don't need to install *nsq* and *Consul*. They are bundled within their Docker images. Also, we don't need *goreman* to run processes.

### Running

```
# bulid all linux binaries
./build.rb binary_all

# build all Docker images
./build.rb image_all

cd datacenters/aws/node1
docker-compose up
```

