# MDTS

## Init

clone source:

```sh
$ git clone https://github.com/MephistoMMM/mdts.git
$ cd mdts
$ export PROJECT_ROOT=$PWD
```

generate protobuf code:

```sh
# activate to project
$ . $PROJECT_ROOT/activate
# generate protobuf code
$ $PROJECT_ROOT/scripts/gen_proto.sh
```

install dependencs:

```sh
# install dep
$ go get -u github.com/golang/dep/cmd/dep
$ cd $PROJECT_ROOT/src/mdts && dep ensure
$ cd $PROJECT_ROOT/src/players && dep ensure
```

## Create Docker Image

```sh
$ cd $PROJECT_ROOT/src/mdts
$ ./build.sh <appname> <version>
```
