#!/usr/bin/env bash
# a hack to generate releases like other prometheus projects
# use like this: 
#       VERSION=1.0.1 ./release.sh


make release_linux
rm -rf "bin/pushprom-$VERSION.linux-amd64"
mkdir "bin/pushprom-$VERSION.linux-amd64"
cp bin/pushprom "bin/pushprom-$VERSION.linux-amd64/pushprom"
cd bin
tar -zcvf "pushprom-$VERSION.linux-amd64.tar.gz" "pushprom-$VERSION.linux-amd64"
