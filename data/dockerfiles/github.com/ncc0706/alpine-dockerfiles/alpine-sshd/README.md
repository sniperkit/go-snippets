docker build -t alpine-sshd .

docker run -d -p 1021:22 --name asd alpine-sshd

