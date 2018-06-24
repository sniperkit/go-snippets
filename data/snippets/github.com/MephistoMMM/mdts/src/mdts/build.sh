#!/bin/bash
# Usage:
#   ./build.sh <app_name>
DOCKER_USER=mpsss
PRO_ROOT=${PWD}

if [ "$1" == "" ] || [ "$2" == "" ]; then
    echo "Error: Lack Argument <app_dir> <app_name>"
    exit 1
fi
PLAYER_DIR=$1
PLAYER_NAME=$2

if [ "$3" == "" ]; then
    VERSION="latest"
else
    VERSION=$3
fi

cd $PLAYER_DIR
go build -o app
cd $PRO_ROOT

docker build --build-arg PLAYER=$PLAYER_DIR -t $PLAYER_NAME:$VERSION .
docker tag $PLAYER_NAME:$VERSION $DOCKER_USER/$PLAYER_NAME:$VERSION
