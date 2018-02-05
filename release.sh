#!/usr/bin/env bash
# a hack to generate releases like other prometheus projects
# use like this: 
#       VERSION=1.0.1 ./release.sh

mkdir -p bin
rm -rf "bin/pushprom-$VERSION.linux-amd64"
mkdir "bin/pushprom-$VERSION.linux-amd64"
env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o "bin/pushprom-$VERSION.linux-amd64/pushprom" github.com/messagebird/pushprom
cd bin
tar -zcvf "pushprom-$VERSION.linux-amd64.tar.gz" "pushprom-$VERSION.linux-amd64"
