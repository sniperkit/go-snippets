set -ex
# docker hub username
USERNAME=$1
VERSION=$2
echo $1
echo $2
# image name
IMAGE=opencv
docker build --no-cache -t $USERNAME/$IMAGE:latest -t $USERNAME/$IMAGE:$VERSION $VERSION
#docker build -t $USERNAME/$IMAGE:latest -t $USERNAME/$IMAGE:$VERSION $VERSION
echo "# $IMAGE $VERSION" >$VERSION/README.md
echo "\`\`\`" >>$VERSION/README.md
docker run --rm $USERNAME/$IMAGE >>$VERSION/README.md
echo "\`\`\`" >>$VERSION/README.md
