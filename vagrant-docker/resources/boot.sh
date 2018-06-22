# docker build -t spring_sample:0.1 .
SERVER_IP="$(ip addr show eth1 | grep 'inet ' | cut -f2 | awk '{ print $2}' | cut -d/ -f1)"
# docker run --privileged -p 0.0.0.0:8080:8080 -p 0.0.0.0:9010:9010 -e SERVER_IP=$SERVER_IP spring_sample:0.1
docker run -p 0.0.0.0:8080:8080 -p 0.0.0.0:9010:9010 -e SERVER_IP=$SERVER_IP spring_sample:0.1

