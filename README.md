# [Redis Alpine Docker Images](https://hub.docker.com/r/comodal/alpine-redis/) [![](https://images.microbadger.com/badges/image/comodal/alpine-redis.svg)](https://microbadger.com/images/comodal/alpine-redis "microbadger.com") [![](https://images.microbadger.com/badges/commit/comodal/alpine-redis.svg)](https://microbadger.com/images/comodal/alpine-redis "microbadger.com")

## Docker Run

```shell
docker run -d\
 --name redis-6379\
 -v /host/dir:/redis/data\
 -p 6379:6379/tcp\
  comodal/alpine-redis:unstable\
   --port 6379\
   --protected-mode no\
   --logfile redis-6379.log
```

## Docker Compose

```yaml
version: '2'

services:
 redis-6379:
  ports:
   - "6379:6379"
  volumes:
   - /host/dir/redis/modules:/redis/modules
   - /host/dir/redis/data:/redis/data
  image: comodal/alpine-redis:unstable
  command: ['--port', '6379', '--protected-mode', 'no', 'logfile', 'redis-6379.log']
```
