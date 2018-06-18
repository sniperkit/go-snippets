Build the app as a static binary and test it's functionality.

```console
cd io16/app/monolith
go build -tags netgo -ldflags "-extldflags '-lm -lstdc++ -static'" .
./monolith --http :10180 --health :10181 &
curl http://127.0.0.1:10180
curl http://127.0.0.1:10180/secure
curl http://127.0.0.1:10180/login -u user
curl http://127.0.0.1:10180/secure -H "Authorization: Bearer <token>"
```

Once we have a binary, we can use Docker to package and distribute it.

```console
docker build -t askcarter/monolith:1.0.0 .
docker push askcarter/monolith:1.0.0
docker run -d askcarter/monolith:1.0.0
docker ps
docker inspect <cid>
curl http://<docker-ip>
docker rm <cid>
docker rmi askcarter/monolith:1.0.0
```
