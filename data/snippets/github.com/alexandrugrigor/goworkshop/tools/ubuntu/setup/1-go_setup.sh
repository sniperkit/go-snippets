#!/bin/bash

export GO_TAR="go1.9.4.linux-amd64.tar.gz"

printf "\n*************************\n"
printf "* 1. Downloading golang *\n"
printf "*************************\n\n"

wget https://dl.google.com/go/$GO_TAR

# check download
if [ -e $GO_TAR ]
then
    echo "Golang successfully downloaded."
else
    echo "Failed to download golang. Check logs for more details."
    exit 1
fi


printf "\n*********************************************************\n"
printf "* 2. Extracting golang (may require your user password) *\n"
printf "*********************************************************\n\n"
sudo tar -C /usr/local -xzf $GO_TAR

echo "checking golang instalation"
if [ -e /usr/local/go ]
then
    echo "Golang successfully installed."
else
    echo "Failed to install golang."
    exit 1
fi

echo "Adding golang to your PATH variable"
printf "\nexport PATH=\$PATH:/usr/local/go/bin\n" >> $HOME/.profile
source $HOME/.profile
echo "Path after adding golang: PATH=$PATH"

printf "\n*********************************************************\n"
printf "* 3.Testing golang instalation                          *\n"
printf "*********************************************************\n\n"

echo "creating the test workspace"
mkdir -p testWorkspace/src/hello
echo "creating the hello.go file"
touch testWorkspace/src/hello/hello.go
echo "adding the hello world code"
printf "package main\n\nimport \"fmt\"\n\nfunc main() {\n    fmt.Printf(\"hello world!\")\n}\n" >> testWorkspace/src/hello/hello.go
export GOPATH=$PWD/testWorkspace
go build hello
export hello_output=$(./hello)
echo "Program output:$hello_output"

if [ "$hello_output" == "hello world!" ]; then
  echo "Go successfully set"
else
  echo "A program occurred while setting up go. Check logs for possible problems"
  exit 1
fi
#cleanup
rm -rf testWorkspace/
rm hello
rm $GO_TAR
