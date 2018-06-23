# tensorflow-go docker installation

git clone this repo

```
docker build -t tensorflow-go .
docker run -p 8888:8888 -d tensorflow-go
docker ps

CONTAINER ID        IMAGE               COMMAND             CREATED             STATUS              PORTS                              NAMES
e726f3ee010c        tensorflow-go       "/run_jupyter.sh"   25 seconds ago      Up 23 seconds       6006/tcp, 0.0.0.0:8888->8888

docker exec -it <CONTAINER ID> bash
```

and your in.. you now have tensorflow + golang availble. now `go get` your favorite open source go for data science:

[gopherdata](https://github.com/gopherdata/resources)
