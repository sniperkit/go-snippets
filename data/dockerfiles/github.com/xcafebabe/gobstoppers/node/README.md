# Alpine Node Images

## Node Image

In order to build your own Node Image for your Raspberry Pi based on `gliderlabs/alpine`

```
git clone https://github.com/xcafebabe/gobstoppers.git && \
cd gobstoppers/node && \
sudo docker build -t kafebob/rpi-alpine-node:latest .
```

Be patient is going to take a while.

If you just want to use latest image available, pull a copy from DockerHub

```
sudo docker pull kafebob/rpi-alpine-node:latest
```
