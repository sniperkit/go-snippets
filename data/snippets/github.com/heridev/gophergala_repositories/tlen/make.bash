#!/usr/bin/env bash

# Copyright 2014 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

set -e

if [ ! -f make.bash ]; then
	exit 1
fi

mkdir -p libs/armeabi-v7a src/go/android
ANDROID_APP=$PWD

rm -rf gen bin

ln -sf $GOPATH/src/golang.org/x/mobile/app/*.java $ANDROID_APP/src/go
ln -sf $GOPATH/src/golang.org/x/mobile/bind/java/Seq.java $ANDROID_APP/src/go

mkdir -p $ANDROID_APP/src/com/
mkdir -p $ANDROID_APP/src/go/rpc/
ln -sf $GOPATH/src/github.com/seletskiy/go-android-rpc/src/com/goandroidrpc/ $ANDROID_APP/src/com/
ln -sf $GOPATH/src/github.com/seletskiy/go-android-rpc/src/go/rpc/ $ANDROID_APP/src/go/rpc/

#./generate.bash
CGO_ENABLED=1 GOOS=android GOARCH=arm GOARM=7 \
	go build -ldflags="-shared" .
mv -f tlen libs/armeabi-v7a/libgojni.so
ant debug
