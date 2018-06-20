# What is this repo for?
This repo contains an example of using TensorFlow Go's API for image recognition
on a pre-trained [model](https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip). 
The model is of the [Inception Architecture](http://arxiv.org/abs/1512.00567).

# Installation
[Instructions for Tensorflow Go Installation](https://www.tensorflow.org/install/install_go)
[Simplified Instructions for Tensorflow Go Installation using Docker](https://github.com/ctava/tensorflow-go)

# Getting it
go get github.com/ctava/tensorflow-go-imagerecognition

# Running it
go run main.go -dir=./ -image=./grumpycat.jpg

# Results
![IMAGE](./grumpycat.jpg)

BEST MATCH: (66% likely) Persian cat
