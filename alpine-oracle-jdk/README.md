# alpine-redis

docker build -t alpine-redis .

docker run -d -p 1908:22 -p 9000:6379  --name redis alpine-redis