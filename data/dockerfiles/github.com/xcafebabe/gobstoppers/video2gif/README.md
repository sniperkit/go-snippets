# Video to Gif

A tool to convert videos (wmv, mp4, mkv and much more) into gif files.

## How to use it

Assuming your video is called `video.mp4` and you want to transform it into a gif called `video.gif`, you will have to run:

```
docker run --rm -v=`pwd`:/tmp/ffmpeg kafebob/video2gif video.mp4 video.gif
```

Container will also accept these options

* `-f <framerate>` 10 fps is the default and works well, but in many cases you can even get good results with lower framerates.
* `-y <yscale>` 480 px wide is the default and automatically determine the height while preserving the aspect ratio.
* `-s <skip_sec>` 0 is default and will skip the first N seconds of video.
* `-t <duration_sec>` If specified creates a N second output.

See how it looks previous example command using some options:

```
docker run --rm -v=`pwd`:/tmp/ffmpeg kafebob/video2gif -f 30 -y 320 -s 10 -t 5 video.mp4 video.gif
```

It will generate a gif with the following specifications:

* 30 frames per second.
* 320 pixels wide preserving the aspect ratio.
* Starts at second 10.
* 5 seconds long.

Read my [blog post](https://toub.es/2017/09/11/high-quality-gif-with-ffmpeg-and-docker/) about video2gif image.
