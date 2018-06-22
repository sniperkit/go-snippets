# sshd

docker build -t alpine-sshd-test .

docker run -d -p 1021:22 --name asd alpine-sshd-test

wget 支持https

