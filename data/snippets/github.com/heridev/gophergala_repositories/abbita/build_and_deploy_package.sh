#! /bin/bash
set -e
ginkgo -r --randomizeAllSpecs -cover
bower install
OUTPUT_FOLDER='builds'
CURRENT_REVISION=`git rev-parse --short  HEAD`
OUTPUT_FILE="$OUTPUT_FOLDER/$CURRENT_REVISION.tar.gz"
echo "Building application"
APP_NAME="abbita"
go get 
GOOS=linux GOARCH=amd64 go build -o $APP_NAME
# docker run --rm -v "$(pwd)":/go/src/github.com/gophergala/abbita -w /go/src/github.com/gophergala/abbita  -e GOOS=linux -e GOARCH=amd64 golang:1.3 go get && go build -o $APP_NAME
tar -zcf $OUTPUT_FILE $APP_NAME public
rm -f $APP_NAME
echo "Application saved to $OUTPUT_FILE"
cp $OUTPUT_FILE builds/latest.tar.gz
echo "Deploying app"
ansible-playbook -i provision/inv.ini provision/app.yml 