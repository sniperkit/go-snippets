# docker-openalpr-alpine

OpenALPR (https://github.com/openalpr/openalpr) running in Docker image on Alpine linux.

## Usage

```bash
docker run --rm -ti -v $PWD:/openalpr wallneradam/docker-openalpr-alpine -c eu test.jpg
```

### Calibrate

You can use the calibrate utility as well on Linux under X:

```bash
docker run --rm -ti \
    --entrypoint=openalpr-utils-calibrate \
    -e DISPLAY -v /tmp:/tmp --ipc=host \
    -v $PWD:/openalpr wallneradam/docker-openalpr-alpine -c eu test.jpg
```

The line started with "-e DISPLAY..." makes container be able to use your X session and able to show the gui.

