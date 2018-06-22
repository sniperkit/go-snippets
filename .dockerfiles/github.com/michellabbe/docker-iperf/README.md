[![Docker Stars](https://img.shields.io/docker/stars/mlabbe/iperf.svg)](https://hub.docker.com/r/mlabbe/iperf/) [![Docker Pulls](https://img.shields.io/docker/pulls/mlabbe/iperf.svg)](https://hub.docker.com/r/mlabbe/iperf/)

# docker: iperf

This is a Docker image to run the [iperf](http://sourceforge.net/projects/iperf/) commandline tool in server mode.
iPerf (iPerf v2) is a tool for active measurements of the maximum achievable bandwidth on IP networks.
It supports tuning of various parameters related to timing, buffers and protocols (TCP, UDP, with IPv4 and IPv6).
For each test it reports the bandwidth and other parameters.

________________________________________
```sh
C:\Users\Example>iperf -c speedtest.mydomain.local
------------------------------------------------------------
Client connecting to speedtest.mydomain.local, TCP port 5001
TCP window size:  208 KByte (default)
------------------------------------------------------------
[  3] local 192.168.1.100 port 53912 connected with 192.168.1.50 port 5001
[ ID] Interval       Transfer     Bandwidth
[  3]  0.0-10.0 sec  1.10 GBytes   941 Mbits/sec
```
________________________________________
```sh
C:\Users\Example>iperf -u -b 50M -c speedtest.mydomain.local
------------------------------------------------------------
Client connecting to speedtest.mydomain.local, UDP port 5001
Sending 1470 byte datagrams, IPG target: 235.20 us (kalman adjust)
UDP buffer size:  208 KByte (default)
------------------------------------------------------------
[  3] local 192.168.1.100 port 52186 connected with 192.168.1.50 port 5001
[ ID] Interval       Transfer     Bandwidth
[  3]  0.0-10.0 sec  59.6 MBytes  49.9 Mbits/sec
[  3] Sent 42501 datagrams
[  3] WARNING: did not receive ack of last datagram after 10 tries.
```
________________________________________
- Runs as non-root user
- Small image size
- Small number of layers
- Supports TCP and UDP

Total size of this image is only:

[![](https://images.microbadger.com/badges/image/mlabbe/iperf.svg)](https://microbadger.com/images/mlabbe/iperf)

________________________________________
### Pulling from Docker hub
If you want to obtain the image from Docker registry, you can use the following command:
```sh
docker pull mlabbe/iperf
```
________________________________________
### Running the image
In order to run the iperf server, use the following:
```sh
docker run --restart=unless-stopped --name=iperf -d -p 5001:5001/tcp -p 5001:5001/udp mlabbe/iperf
```
At that point, you can use your Docker server as an iperf server to begin testing your network.