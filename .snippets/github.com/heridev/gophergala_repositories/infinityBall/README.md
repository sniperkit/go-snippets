# infinityBall

# Idea

The idea was to have a application to run a randomly generated terrain for
a physical object to interact with. The application can be run indefinite
as it generates terrain for the object if the object is close to the
boundaries.

The first implementation is a ball that is able to interact with the
environment.

# How to use 

To use:

Linux distribution using go1.4.1

export GOBIN=$GOPATH/bin 

export PATH=$GOPATH/bin:$PATH:/usr/local/go/bin

comment can be added to /.bashrc

sudo apt-get update
sudo apt-get install libglew-dev libxrandr-dev libglfw-dev

go get github.com/go-gl/gl
go get github.com/go-gl/glu
go get github.com/go-gl/mathgl/mgl64
go get github.com/go-gl/glh
go get github.com/go-gl/glfw

go install infinityBall/;

infinityBall

# TODO

We had a hard time using the exciting physics engines available for
Golang. A migration to Unity would make it easier for rendering objects as
well as adding physics.
