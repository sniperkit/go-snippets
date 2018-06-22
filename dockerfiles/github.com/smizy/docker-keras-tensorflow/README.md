# docker-keras-tensorflow
[![](https://images.microbadger.com/badges/image/smizy/keras.svg)](https://microbadger.com/images/smizy/keras "Get your own image badge on microbadger.com") 
[![](https://images.microbadger.com/badges/version/smizy/keras.svg)](https://microbadger.com/images/smizy/keras "Get your own version badge on microbadger.com")
[![build status](https://gitlab.com/smizy/docker-keras-tensorflow/badges/master/build.svg)](https://gitlab.com/smizy/docker-keras-tensorflow/commits/master)

Python3 Keras(Tensorflow backended) with Jupyter docker image based on alpine 

* numpy, scipy, pandas, scikit-learn, seaborn, tensorflow, keras installed via pip. See `pip list --format=columns` for detail.
* CPU only

## Usage
```
docker run -it --rm -v $(pwd):/data -w /data -p 8888:8888 smizy/keras:2.0-cpu-alpine
```