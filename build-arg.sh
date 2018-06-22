set -ex
USERNAME=$1	# docker hub username
APP_VER_BASEOS_VER=$2	# 3.3.0-alpine -> 3.3.0-alpine:3.5 "opencv:3.3.0-alpine:3.5"
echo "1=$1"
echo "2=$2"
APP_VERSION=$(echo $APP_VER_BASEOS_VER | cut -f1 -d-)
APP_VERSION=$(echo $APP_VERSION | cut -f2 -d:)
BASEOS_VER=$(echo $APP_VER_BASEOS_VER | cut -f2 -d-)
OS=$(echo $BASEOS_VER | cut -f1 -d:)
IMAGE=$(echo $APP_VER_BASEOS_VER | cut -f1 -d:)
echo "USERNAME=$USERNAME"
echo "APP_VERSION=$APP_VERSION"
echo "BASEOS_VER=$BASEOS_VER"
echo "IMAGE=$IMAGE"
IMAGE=opencv # image name
LOCAL_DIR=$APP_VERSION-$OS
docker build --no-cache --build-arg BASEOS_VER=$BASEOS_VER --build-arg APP_VER=$APP_VERSION -t $USERNAME/$IMAGE:latest -t $USERNAME/$IMAGE:$APP_VERSION $LOCAL_DIR
echo "# $IMAGE $APP_VERSION $BASE" >$LOCAL_DIR/README.md
echo "\`\`\`" >>$LOCAL_DIR/README.md
docker run --rm $USERNAME/$IMAGE >>$LOCAL_DIR/README.md
echo "\`\`\`" >>$LOCAL_DIR/README.md
